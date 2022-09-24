package dao

import "time"

type Analyses struct {
	Id        int
	Url       string
	Requested time.Time
	Completed *time.Time
}
