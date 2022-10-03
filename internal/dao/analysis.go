package dao

import "time"

type Analyses struct {
	Id                int64
	Url               string
	Requested         time.Time
	Completed         *time.Time
	Title             string
	Headings          map[string]int
	InternalLinkCount int
	ExternalLinkCount int
	ActiveLinkCount   int
	InactiveLinkCount int
	PageVersion       string
	HasLoginForm      bool
}
