package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DiLRandI/web-analyser/internal/dao"
	"github.com/DiLRandI/web-analyser/internal/dto"
	"github.com/DiLRandI/web-analyser/internal/repository"
	"github.com/DiLRandI/web-analyser/internal/repository/mem"
	"github.com/DiLRandI/web-analyser/internal/service/webpage"
	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	"github.com/sirupsen/logrus"
)

type Processor interface {
	ProcessPage(ctx context.Context, req *dto.AnalysesRequest) (*dto.AnalysesResponse, error)
	GetProcessResultFor(ctx context.Context, id int64) (*dto.ResultResponse, error)
	GetProcessResults(ctx context.Context) ([]*dto.ResultResponse, error)
}

type processor struct {
	downloader webpage.Downloader
	analyserFn func() webpage.Analyser
	result     repository.Results
}

func NewProcessor(downloader webpage.Downloader,
	analyserFn func() webpage.Analyser,
	result repository.Results) Processor {
	return &processor{
		downloader: downloader,
		analyserFn: analyserFn,
		result:     result,
	}
}

func (s *processor) ProcessPage(
	ctx context.Context, req *dto.AnalysesRequest,
) (*dto.AnalysesResponse, error) {
	if req.WebUrl == "" {
		return nil, fmt.Errorf("WebUrl is required.")
	}

	m, err := s.downloader.Download(ctx, req.WebUrl)
	if err != nil {
		return nil, err
	}

	if m.Content == nil {
		return nil, fmt.Errorf("there is no content to process further")
	}

	id, err := s.result.Save(ctx, &dao.Analyses{
		Url:           m.Url,
		Requested:     time.Now(),
		ProcessStatus: &dao.ProcessStatusCreated,
	})
	if err != nil {
		return nil, err
	}

	go s.bgProcess(id, m)

	return &dto.AnalysesResponse{Id: id}, nil
}

func (s *processor) GetProcessResults(ctx context.Context) ([]*dto.ResultResponse, error) {
	results, err := s.result.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := []*dto.ResultResponse{}
	for _, r := range results {
		m := &dto.ResultResponse{}
		m.Id = r.Id
		m.Url = r.Url
		m.Requested = r.Requested
		m.Completed = r.Completed
		m.ProcessStatus = string(*r.ProcessStatus)
		m.Title = r.Title
		m.Headings = r.Headings
		m.InternalLinkCount = r.InternalLinkCount
		m.ExternalLinkCount = r.ExternalLinkCount
		m.ActiveLinkCount = r.ActiveLinkCount
		m.InactiveLinkCount = r.InactiveLinkCount
		m.PageVersion = r.PageVersion
		m.HasLoginForm = r.HasLoginForm

		res = append(res, m)
	}

	return res, nil
}

func (s *processor) GetProcessResultFor(ctx context.Context, id int64) (*dto.ResultResponse, error) {
	result, err := s.result.Get(ctx, id)
	if err != nil {
		if errors.Is(err, mem.ResultNotFoundErr) {
			return nil, &NotFoundError{msg: err.Error()}
		}

		return nil, err
	}

	res := &dto.ResultResponse{}
	res.Id = result.Id
	res.Url = result.Url
	res.Requested = result.Requested
	res.Completed = result.Completed
	res.ProcessStatus = string(*result.ProcessStatus)
	res.Title = result.Title
	res.Headings = result.Headings
	res.InternalLinkCount = result.InternalLinkCount
	res.ExternalLinkCount = result.ExternalLinkCount
	res.ActiveLinkCount = result.ActiveLinkCount
	res.InactiveLinkCount = result.InactiveLinkCount
	res.PageVersion = result.PageVersion
	res.HasLoginForm = result.HasLoginForm

	return res, nil
}

func (s *processor) bgProcess(id int64, m *model.DownloadedWebpage) {
	logrus.Infof("Starting analysis for url %q", m.Url)
	ctx := context.Background()
	analysis, err := s.result.Get(ctx, id)
	if err != nil {
		logrus.Error(err)
		return
	}

	svc := s.analyserFn()
	pageResult, err := svc.AnalysePage(ctx, m)
	if err != nil {
		logrus.Error(err)
		s.updateProcessStatus(ctx, id, analysis, dao.ProcessStatusFailed)
		return
	}

	logrus.Infof("Analysis completed, %+#v", pageResult)
	analysis.Completed = timePtr(time.Now())
	analysis.ProcessStatus = &dao.ProcessStatusCompleted
	analysis.Title = pageResult.Title
	analysis.Headings = pageResult.Headings
	analysis.InternalLinkCount = pageResult.InternalLinkCount
	analysis.ExternalLinkCount = pageResult.ExternalLinkCount
	analysis.ActiveLinkCount = pageResult.ActiveLinkCount
	analysis.InactiveLinkCount = pageResult.InactiveLinkCount
	analysis.PageVersion = pageResult.PageVersion
	analysis.HasLoginForm = pageResult.HasLoginForm

	logrus.Infof("updating the result, %+#v", analysis)
	if err := s.result.Update(ctx, id, analysis); err != nil {
		logrus.Errorf("unable to update the results, %v", err)
	}
}

func (s *processor) updateProcessStatus(
	ctx context.Context, id int64, m *dao.Analyses, ps dao.ProcessStatus,
) {
	m.ProcessStatus = &ps
	if err := s.result.Update(ctx, id, m); err != nil {
		logrus.Errorf("updating process status to %q failed for analysis id %d", ps, id)
	}

}

func timePtr(t time.Time) *time.Time {
	return &t
}
