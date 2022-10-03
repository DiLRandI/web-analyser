package dao

import "time"

type Analyses struct {
	Id                int64
	Url               string
	Requested         time.Time
	Completed         *time.Time
	ProcessStatus     *ProcessStatus
	Title             string
	Headings          map[string]int
	InternalLinkCount int
	ExternalLinkCount int
	ActiveLinkCount   int
	InactiveLinkCount int
	PageVersion       string
	HasLoginForm      bool
}

type ProcessStatus string

var (
	ProcessStatusCreated   ProcessStatus = "Created"
	ProcessStatusCompleted ProcessStatus = "Completed"
	ProcessStatusFailed    ProcessStatus = "Failed"
)
