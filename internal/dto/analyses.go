package dto

type AnalysesRequest struct {
	WebUrl string `json:"webUrl"`
}

type AnalysesResponse struct {
	Id int64 `json:"id"`
}

type ResultRequest struct {
	Id int64 `json:"id"`
}
