package webpage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	"golang.org/x/net/html"
)

type Analyser interface {
}

type analyser struct {
}

func NewAnalyser() Analyser {
	return &analyser{}
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

func (s *analyser) linksDetail(ctx context.Context, node *html.Node) (any, error) {
	return nil, nil
}

func (s *analyser) hasLoginForm(ctx context.Context, node *html.Node) (bool, error) {
	return false, nil
}
