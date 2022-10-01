package model

type DownloadedWebpage struct {
	StatusCode int
	Status     string
	Url        string
	Content    []byte
}

type Analysis struct {
}
