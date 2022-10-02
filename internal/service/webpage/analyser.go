package webpage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Analyser interface {
}

type analyser struct {
	client WebClient
}

func NewAnalyser(client WebClient) Analyser {
	return &analyser{
		client: client,
	}
}

func (s *analyser) AnalysePage(ctx context.Context, page *model.DownloadedWebpage) (*model.Analysis, error) {

	return nil, nil
}

func (s *analyser) pageVersion(ctx context.Context, content []byte) (string, error) {
	tt := html.NewTokenizer(bytes.NewReader(content))
	docType := ""
loop:
	for {
		token := tt.Next()
		switch token {
		case html.ErrorToken:
			err := tt.Err()
			if errors.Is(err, io.EOF) {
				return "", fmt.Errorf("!Doctype node is not found in the document")
			}

			return "", fmt.Errorf("Unable to process the document, %v", err)

		case html.DoctypeToken:
			docType = string(tt.Text())
			break loop
		}
	}

	//html5
	if docType == "html" {
		return "HTML5 and beyond", nil
	}

	// old / other doc types
	dts := strings.Split(docType, "-//")
	if len(dts) != 2 {
		return "", fmt.Errorf("Unable to parse the Doctype node %q", docType)
	}

	dts = strings.Split(dts[1], "//")
	if len(dts) > 2 {
		return dts[1], nil
	}

	return "", fmt.Errorf("Unable to parse the Doctype node %q", docType)
}

func (s *analyser) pageTitle(ctx context.Context, content []byte) (string, error) {
	tt := html.NewTokenizer(bytes.NewReader(content))
	title := ""
	foundTitleNode := false
	foundHeadElement := false
loop:
	for {
		token := tt.Next()
		switch token {
		case html.ErrorToken:
			err := tt.Err()
			if errors.Is(err, io.EOF) {
				return "", fmt.Errorf("head element not found in the document")
			}

			return "", fmt.Errorf("unable to process the document, %v", err)

		case html.StartTagToken:
			name, _ := tt.TagName()
			if string(name) == "head" {
				foundHeadElement = true
			}
			if foundHeadElement && string(name) == "title" {
				foundTitleNode = true
			}
		case html.EndTagToken:
			name, _ := tt.TagName()
			if string(name) == "head" {
				break loop
			}
		case html.TextToken:
			if foundTitleNode {
				title = strings.TrimSpace(string(tt.Text()))
				break loop
			}
		}
	}

	if !foundTitleNode {
		return "", fmt.Errorf("title node not found in the document")
	}
	return title, nil
}

func (s *analyser) headingDetails(ctx context.Context, content []byte) (map[string]int, error) {
	headings := map[string]int{
		"h1": 0,
		"h2": 0,
		"h3": 0,
		"h4": 0,
		"h5": 0,
		"h6": 0,
	}

	tt := html.NewTokenizer(bytes.NewReader(content))
	for {
		token := tt.Next()
		switch token {
		case html.ErrorToken:
			err := tt.Err()
			if errors.Is(err, io.EOF) {
				return headings, nil
			}

			return nil, fmt.Errorf("unable to process the document, %v", err)

		case html.StartTagToken:
			name, _ := tt.TagName()
			switch string(name) {
			case "h1":
				headings["h1"]++
			case "h2":
				headings["h2"]++
			case "h3":
				headings["h3"]++
			case "h4":
				headings["h4"]++
			case "h5":
				headings["h5"]++
			case "h6":
				headings["h6"]++
			}
		}
	}
}

// hasLoginForm to detect login form first check to see if there is a <form> element,
// if <form> element found then it will check for <input type="password"> to make sure
// that the password is asked, finally check for either <button type="submit">Login</button> or
// <input type="submit" value="Login">
// Limitation only check for "Login" text to make it simpler,
func (s *analyser) hasLoginForm(ctx context.Context, content []byte) (bool, error) {
	insideForm := false
	hasPasswordInput := false
	hasLoginButton := false
	possibleLoginButton := false

	tt := html.NewTokenizer(bytes.NewReader(content))
	for {
		token := tt.Next()
		switch token {
		case html.ErrorToken:
			err := tt.Err()
			if errors.Is(err, io.EOF) {
				return false, nil
			}

			return false, fmt.Errorf("unable to process the document, %v", err)

		case html.StartTagToken:
			name, hasAttr := tt.TagName()
			if string(name) == "form" {
				insideForm = true
				continue
			} else if insideForm && string(name) == "input" && hasAttr {
				isSubmit := false
				isLogin := false
				for {
					k, v, m := tt.TagAttr()
					sk := string(k)
					sv := string(v)
					if sk == "type" && sv == "password" {
						hasPasswordInput = true
						break
					}
					if sk == "type" && sv == "submit" {
						isSubmit = true
					}
					if sk == "value" && sv == "Login" {
						isLogin = true
					}

					if !m {
						break
					}
				}

				hasLoginButton = isSubmit && isLogin
			} else if insideForm && string(name) == "button" && hasAttr {
				// can be possible login button
				for {
					k, v, m := tt.TagAttr()
					sk := string(k)
					sv := string(v)
					if sk == "type" && sv == "submit" {
						possibleLoginButton = true
						break
					}

					if !m {
						break
					}
				}

			}
		case html.EndTagToken:
			name, _ := tt.TagName()
			if string(name) == "form" {
				// is closing form is a login form ?
				if hasLoginButton && hasPasswordInput {
					return true, nil
				}
				insideForm = false
				hasPasswordInput = false
				hasLoginButton = false
				possibleLoginButton = false
			} else if string(name) == "button" {
				possibleLoginButton = false
			}
		case html.TextToken:
			if possibleLoginButton && !hasLoginButton {
				hasLoginButton = string(tt.Text()) == "Login"
			}
		}
	}
}

func (s *analyser) linksDetail(ctx context.Context, hostUrl string, content []byte) ([]*model.Link, error) {
	links := []*model.Link{}
	insideLinkTak := false
	tt := html.NewTokenizer(bytes.NewReader(content))
	wg := sync.WaitGroup{}
	for {
		token := tt.Next()
		switch token {
		case html.ErrorToken:
			err := tt.Err()
			if errors.Is(err, io.EOF) {
				wg.Wait()
				return links, nil
			}

			return nil, fmt.Errorf("unable to process the document, %v", err)
		case html.StartTagToken, html.EndTagToken:
			name, atr := tt.TagName()
			if string(name) != "a" {
				continue
			}

			if !insideLinkTak {
				links = append(links, &model.Link{})
			}

			insideLinkTak = !insideLinkTak
			if !insideLinkTak || !atr {
				continue
			}

			for {
				k, v, m := tt.TagAttr()
				sk := string(k)
				sv := string(v)
				if sk == "href" {
					links[len(links)-1].Url = sv
					wg.Add(1)
					go func(l *model.Link, host, link string) {
						defer wg.Done()
						l.IsInternal = s.isInternalLink(hostUrl, sv)
						status, code := s.linkStatus(hostUrl, sv)
						l.LinkStatus = status
						l.HttpStatusCode = code
					}(links[len(links)-1], hostUrl, sv)

				}

				if !m {
					break
				}
			}
		case html.TextToken:
			if insideLinkTak {
				links[len(links)-1].Name = string(tt.Text())
			}
		}

	}
}

func (s *analyser) isInternalLink(host, link string) bool {
	hostUrl, err := url.ParseRequestURI(host)
	if err != nil {
		logrus.Errorf("unable to parse hostUrl %q, %v", host, err)
		return false
	}

	linkUrl, err := url.Parse(link)
	if err != nil {
		logrus.Errorf("unable to parse linkUrl %q, %v", link, err)
		return false
	}

	if linkUrl.Hostname() == "" ||
		linkUrl.Hostname() == hostUrl.Hostname() {
		return true
	}

	return false
}

func (s *analyser) linkStatus(host, link string) (model.LinkStatus, int) {
	linkUrl, err := url.Parse(link)
	if err != nil {
		logrus.Errorf("unable to parse linkUrl %q, %v", link, err)
		return model.LinkStatusInactive, -1
	}

	getUrl := link
	if linkUrl.Host == "" {
		getUrl = fmt.Sprintf("%s/%s", host, link)
	}

	res, err := s.client.Get(getUrl)
	if err != nil {
		return model.LinkStatusInactive, -1
	}

	if res.StatusCode == http.StatusOK {
		return model.LinkStatusActive, res.StatusCode
	}

	return model.LinkStatusInactive, res.StatusCode
}
