package mock

import (
	"context"

	"github.com/DiLRandI/web-analyser/internal/dto"
	"github.com/stretchr/testify/mock"
)

type ProcessorMock struct {
	mock.Mock
}

func (m *ProcessorMock) ProcessPage(ctx context.Context, req *dto.AnalysesRequest) (*dto.AnalysesResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.AnalysesResponse), args.Error(1)
}
func (m *ProcessorMock) GetProcessResultFor(ctx context.Context, id int64) (*dto.ResultResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.ResultResponse), args.Error(1)
}
func (m *ProcessorMock) GetProcessResults(ctx context.Context) ([]*dto.ResultResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*dto.ResultResponse), args.Error(1)
}
