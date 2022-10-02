package model

type DownloadedWebpage struct {
	StatusCode int
	Status     string
	Url        string
	Content    []byte
}

type Analysis struct {
}

type LinkStatus string

var (
	LinkStatusActive   LinkStatus = "Active"
	LinkStatusInactive LinkStatus = "Inactive"
)

type Link struct {
	Name       string
	Url        string
	IsInternal bool
	LinkStatus LinkStatus
	//HttpStatusCode for inactive links only
	HttpStatusCode int
}
