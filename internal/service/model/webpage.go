package model

type DownloadedWebPage struct {
	StatusCode int
	Status     string
	Url        string
	Content    []byte
}
