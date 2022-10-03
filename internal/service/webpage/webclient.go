package webpage

import (
	"net/http"
)

type WebClient interface {
	Get(url string) (resp *http.Response, err error)
}
