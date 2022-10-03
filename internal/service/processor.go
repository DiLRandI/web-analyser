package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DiLRandI/web-analyser/internal/dao"
	"github.com/DiLRandI/web-analyser/internal/dto"
	"github.com/DiLRandI/web-analyser/internal/repository"
	"github.com/DiLRandI/web-analyser/internal/service/webpage"
	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
)

type Processor interface {
	ProcessPage(ctx context.Context, req *dto.AnalysesRequest) (*dto.AnalysesResponse, error)
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
		Url:       m.Url,
		Requested: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	go s.bgProcess(id, m)

	return &dto.AnalysesResponse{Id: id}, nil
}

func (s *processor) bgProcess(id int64, m *model.DownloadedWebpage) {

	// TODO Bg process
}
