package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DiLRandI/web-analyser/internal/dto"
	"github.com/DiLRandI/web-analyser/internal/service"
	mc "github.com/DiLRandI/web-analyser/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_handler_validations(t *testing.T) {
	//setup
	router := gin.Default()
	sut := &analysisHandler{}
	sut.RegisterRoutes(router)

	testCases := []struct {
		desc          string
		httpMethod    string
		endpoint      string
		payload       io.Reader
		expResponse   string
		expStatusCode int
	}{
		{
			desc:          "analyse handler respond with bad request for invalid body",
			httpMethod:    http.MethodPost,
			endpoint:      "/api/v1/analyse",
			payload:       strings.NewReader(""),
			expStatusCode: http.StatusBadRequest,
		},
		{
			desc:          "analyse handler respond with bad request when WebURL is empty",
			httpMethod:    http.MethodPost,
			endpoint:      "/api/v1/analyse",
			payload:       strings.NewReader(`{}`),
			expStatusCode: http.StatusBadRequest,
		},
		{
			desc:          "getAnalysisById handler respond with bad request when url param id is not valid",
			httpMethod:    http.MethodGet,
			endpoint:      "/api/v1/analyse/abc",
			payload:       nil,
			expStatusCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.httpMethod, tc.endpoint, tc.payload)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expStatusCode, w.Code)

			res, err := io.ReadAll(w.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expResponse, string(res))
		})
	}
}

func Test_handler_service_error(t *testing.T) {
	testCases := []struct {
		desc          string
		httpMethod    string
		endpoint      string
		payload       io.Reader
		expResponse   string
		expStatusCode int
	}{
		{
			desc:          "analyse handler respond with internal server error for service error",
			httpMethod:    http.MethodPost,
			endpoint:      "/api/v1/analyse",
			payload:       strings.NewReader(`{"webUrl":"https://www.test.com/"}`),
			expStatusCode: http.StatusInternalServerError,
		},
		{
			desc:          "getAalyse handler respond with internal server error for service error",
			httpMethod:    http.MethodGet,
			endpoint:      "/api/v1/analyse",
			payload:       nil,
			expStatusCode: http.StatusInternalServerError,
		},
		{
			desc:          "getAnalysisById handler respond internal server error for id 1 generic service error",
			httpMethod:    http.MethodGet,
			endpoint:      "/api/v1/analyse/1",
			payload:       nil,
			expStatusCode: http.StatusInternalServerError,
		},
		{
			desc:          "getAnalysisById handler respond not found for id 2 not found service error",
			httpMethod:    http.MethodGet,
			endpoint:      "/api/v1/analyse/2",
			payload:       nil,
			expStatusCode: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			routeEng := gin.Default()

			mp := new(mc.ProcessorMock)
			mp.On("ProcessPage", mock.Anything, mock.Anything).
				Return((*dto.AnalysesResponse)(nil), errors.New("service failing"))
				//GetProcessResultFor id 1 generic error return
			mp.On("GetProcessResultFor", mock.Anything, int64(1)).
				Return((*dto.ResultResponse)(nil), errors.New("service failing"))
				//GetProcessResultFor id 2 not found error return
			mp.On("GetProcessResultFor", mock.Anything, int64(2)).
				Return((*dto.ResultResponse)(nil), &service.NotFoundError{})
			mp.On("GetProcessResults", mock.Anything).
				Return(([]*dto.ResultResponse)(nil), errors.New("service failing"))

			sut := New(mp)
			sut.RegisterRoutes(routeEng)
			req, _ := http.NewRequest(tc.httpMethod, tc.endpoint, tc.payload)
			routeEng.ServeHTTP(w, req)

			assert.Equal(t, tc.expStatusCode, w.Code)

			res, err := io.ReadAll(w.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expResponse, string(res))
		})
	}
}

func Test_handler_happy_path(t *testing.T) {
	now := time.Now()
	completed := time.Now().Add(time.Second * 5)
	res1 := &dto.ResultResponse{
		Id:            1,
		Url:           "https://www.test.com/",
		Requested:     now,
		Completed:     &completed,
		ProcessStatus: "Completed",
		Title:         "Test",
		Headings: map[string]int{
			"h1": 1,
			"h2": 2,
			"h3": 3,
			"h4": 4,
			"h5": 5,
			"h6": 6,
		},
		InternalLinkCount: 5,
		ExternalLinkCount: 10,
		ActiveLinkCount:   7,
		InactiveLinkCount: 3,
		PageVersion:       "HTML 5",
		HasLoginForm:      true,
	}
	res2 := []*dto.ResultResponse{
		{
			Id:            1,
			Url:           "https://www.test.com/",
			Requested:     now,
			Completed:     &completed,
			ProcessStatus: "Completed",
			Title:         "Test",
			Headings: map[string]int{
				"h1": 1,
				"h2": 2,
				"h3": 3,
				"h4": 4,
				"h5": 5,
				"h6": 6,
			},
			InternalLinkCount: 5,
			ExternalLinkCount: 10,
			ActiveLinkCount:   7,
			InactiveLinkCount: 3,
			PageVersion:       "HTML 5",
			HasLoginForm:      true,
		},
		{
			Id:            2,
			Url:           "https://www.test.com/",
			Requested:     now,
			Completed:     &completed,
			ProcessStatus: "Completed",
			Title:         "Test",
			Headings: map[string]int{
				"h1": 1,
				"h2": 2,
				"h3": 3,
				"h4": 4,
				"h5": 5,
				"h6": 6,
			},
			InternalLinkCount: 5,
			ExternalLinkCount: 10,
			ActiveLinkCount:   7,
			InactiveLinkCount: 3,
			PageVersion:       "HTML 5",
			HasLoginForm:      true,
		},
	}

	res1Json, _ := json.Marshal(res1)
	res2Json, _ := json.Marshal(res2)

	testCases := []struct {
		desc          string
		httpMethod    string
		endpoint      string
		payload       io.Reader
		expResponse   string
		expStatusCode int
	}{
		{
			desc:          "analyse handler 200 with valid id",
			httpMethod:    http.MethodPost,
			endpoint:      "/api/v1/analyse",
			payload:       strings.NewReader(`{"webUrl":"https://www.test.com/"}`),
			expStatusCode: http.StatusAccepted,
			expResponse:   "{\"id\":1}",
		},
		{
			desc:          "getAalyse handler respond with valid response and status 200",
			httpMethod:    http.MethodGet,
			endpoint:      "/api/v1/analyse",
			payload:       nil,
			expStatusCode: http.StatusOK,
			expResponse:   string(res2Json),
		},
		{
			desc:          "getAnalysisById handler respond internal server error for id 1 generic service error",
			httpMethod:    http.MethodGet,
			endpoint:      "/api/v1/analyse/1",
			payload:       nil,
			expStatusCode: http.StatusOK,
			expResponse:   string(res1Json),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			routeEng := gin.Default()
			mp := new(mc.ProcessorMock)
			mp.On("ProcessPage", mock.Anything, mock.Anything).
				Return(&dto.AnalysesResponse{Id: 1}, nil)
				//GetProcessResultFor id 1 generic error return
			mp.On("GetProcessResultFor", mock.Anything, int64(1)).
				Return(res1, nil)
				//GetProcessResultFor id 2 not found error return
			mp.On("GetProcessResults", mock.Anything).
				Return(res2, nil)

			sut := New(mp)
			sut.RegisterRoutes(routeEng)
			req, _ := http.NewRequest(tc.httpMethod, tc.endpoint, tc.payload)
			routeEng.ServeHTTP(w, req)

			assert.Equal(t, tc.expStatusCode, w.Code)

			res, err := io.ReadAll(w.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tc.expResponse), string(res))
		})
	}
}
