package service

import (
	"bytes"
	"context"

	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	"golang.org/x/net/html"
)

type Analyser interface {
}

type analyser struct {
}

func NewAnalyser() Analyser {
	return &analyser{}
}

func (s *analyser) AnalysePage(ctx context.Context, page *model.DownloadedWebpage) (*model.Analysis, error) {
	_, err := html.Parse(bytes.NewReader(page.Content))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *analyser) pageVersion(ctx context.Context) (string, error) {
	return "", nil
}

func (s *analyser) pageTitle(ctx context.Context) (string, error) {
	return "", nil
}

func (s *analyser) linksDetail(ctx context.Context) (any, error) {
	return nil, nil
}

func (s *analyser) hasLoginForm(ctx context.Context) (bool, error) {
	return false, nil
}
