package service

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DiLRandI/web-analyser/internal/service/model"
)

type WebPageDownloader interface {
}

type webPageDownloader struct {
	client *http.Client
}

func NewWebPageDownloader(client *http.Client) WebPageDownloader {
	return &webPageDownloader{
		client: client,
	}
}

func (s *webPageDownloader) Download(url string) (*model.DownloadedWebPage, error) {
	res, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download the webpage, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return &model.DownloadedWebPage{
			StatusCode: res.StatusCode,
			Status:     res.Status,
		}, nil
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read the body content, %v", err)
	}

	return &model.DownloadedWebPage{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Url:        url,
		Content:    content,
	}, nil
}
