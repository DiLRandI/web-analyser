package webpage

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	mc "github.com/DiLRandI/web-analyser/mock"
	"github.com/stretchr/testify/assert"
)

func Test_download(t *testing.T) {
	testCases := []struct {
		desc  string
		url   string
		err   error
		res   *model.DownloadedWebpage
		mcRes *http.Response
		mcErr error
	}{
		{
			desc:  "Download should return an error if webpage can't be downloaded",
			url:   "http://test.com",
			err:   errors.New("unable to download the webpage, test failure"),
			res:   nil,
			mcRes: nil,
			mcErr: errors.New("test failure"),
		},
		{
			desc: "Download should return the status and status code only for 4xx and above statuses",
			url:  "http://test.com",
			err:  nil,
			res: &model.DownloadedWebpage{
				StatusCode: http.StatusNotFound,
				Status:     "404 NOT FOUND",
				Url:        "http://test.com",
			},
			mcRes: &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(strings.NewReader("")),
			},
			mcErr: nil,
		},
		{
			desc: "Download should return the downloaded web page model with valid content",
			url:  "http://test.com",
			err:  nil,
			res: &model.DownloadedWebpage{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Url:        "http://test.com",
				Content:    []byte("test content"),
			},
			mcRes: &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("test content")),
			},
			mcErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mc := new(mc.WebClientMock)
			mc.On("Get", tc.url).Return(tc.mcRes, tc.mcErr)
			sut := NewDownloader(mc)
			res, err := sut.Download(tc.url)

			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}

			if res == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, res, tc.res)
			}
		})
	}
}
