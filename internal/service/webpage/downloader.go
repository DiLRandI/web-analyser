package webpage

import (
	"fmt"
	"io"

	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	"github.com/sirupsen/logrus"
)

type Downloader interface {
	Download(url string) (*model.DownloadedWebpage, error)
}

type downloader struct {
	client WebClient
}

func NewDownloader(client WebClient) Downloader {
	return &downloader{
		client: client,
	}
}

func (s *downloader) Download(url string) (*model.DownloadedWebpage, error) {
	res, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download the webpage, %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Warnf("Error while closing the response body, %v", err)
		}
	}()

	if res.StatusCode >= 400 {
		return &model.DownloadedWebpage{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Url:        url,
		}, nil
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read the body content, %v", err)
	}

	return &model.DownloadedWebpage{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Url:        url,
		Content:    content,
	}, nil
}
