package model

type DownloadedWebpage struct {
	StatusCode int
	Status     string
	Url        string
	Content    []byte
}

type Analysis struct {
	Page              *DownloadedWebpage
	Links             []*Link
	Title             string
	Headings          map[string]int
	InternalLinkCount int
	ExternalLinkCount int
	ActiveLinkCount   int
	InactiveLinkCount int
	PageVersion       string
	HasLoginForm      bool
}

type LinkStatus string

var (
	LinkStatusActive   LinkStatus = "Active"
	LinkStatusInactive LinkStatus = "Inactive"
)

type Link struct {
	Name           string
	Url            string
	IsInternal     bool
	LinkStatus     LinkStatus
	HttpStatusCode int
}
