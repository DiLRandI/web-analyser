package dto

import "time"

type ResultResponse struct {
	Id                int64          `json:"id"`
	Url               string         `json:"url"`
	Requested         time.Time      `json:"requested"`
	Completed         *time.Time     `json:"completed"`
	ProcessStatus     string         `json:"processStatus"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	InternalLinkCount int            `json:"internalLinkCount"`
	ExternalLinkCount int            `json:"externalLinkCount"`
	ActiveLinkCount   int            `json:"activeLinkCount"`
	InactiveLinkCount int            `json:"inactiveLinkCount"`
	PageVersion       string         `json:"pageVersion"`
	HasLoginForm      bool           `json:"hasLoginForm"`
}
