package webpage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_page_version_for_html5(t *testing.T) {
	sut := &analyser{}
	content := `<!DOCTYPE html>`
	version, err := sut.pageVersion(context.Background(), []byte(content))

	assert.NoError(t, err)
	assert.Equal(t, "HTML5 and beyond", version)
}

func Test_page_version_return_an_error_when_doctype_element_is_missing(t *testing.T) {
	sut := &analyser{}
	content := `<html lang="en">
	</html>`
	version, err := sut.pageVersion(context.Background(), []byte(content))

	assert.Empty(t, version)
	assert.ErrorContains(t, err, "!Doctype node is not found in the document")
}

func Test_page_version_for_html_4_01_strict(t *testing.T) {
	sut := &analyser{}
	content := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
	"http://www.w3.org/TR/html4/strict.dtd">`
	version, err := sut.pageVersion(context.Background(), []byte(content))

	assert.NoError(t, err)
	assert.Equal(t, "DTD HTML 4.01", version)
}

func Test_page_version_return_an_error_when_doctype_content_is_malformed(t *testing.T) {
	sut := &analyser{}
	malformedContent := `HTML PUBLIC "//W3C//DTD HTML 4.01//EN"
	"http://www.w3.org/TR/html4/strict.dtd"`
	content := fmt.Sprintf(`<!DOCTYPE %s>`, malformedContent)
	version, err := sut.pageVersion(context.Background(), []byte(content))

	assert.Empty(t, version)
	assert.ErrorContains(t, err, fmt.Sprintf("Unable to parse the Doctype node %q", malformedContent))
}

func Test_page_version_return_an_error_when_doctype_content_version_is_malformed(t *testing.T) {
	sut := &analyser{}
	malformedContent := `HTML PUBLIC "-//W3C"
	"http://www.w3.org/TR/html4/strict.dtd"`
	content := fmt.Sprintf(`<!DOCTYPE %s>`, malformedContent)
	version, err := sut.pageVersion(context.Background(), []byte(content))

	assert.Empty(t, version)
	assert.ErrorContains(t, err, fmt.Sprintf("Unable to parse the Doctype node %q", malformedContent))
}

func Test_page_title_returns_for_valid_title(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">
	<head>
		<title>Test Title</title>
	</head>
	<body>
		
	</body>
	</html>`
	title, err := sut.pageTitle(context.Background(), []byte(content))

	assert.NoError(t, err)
	assert.Equal(t, title, "Test Title")
}

func Test_page_title_returns_empty_for_empty_title_node(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">
	<head>
		<title></title>
	</head>
	<body>
		
	</body>
	</html>`
	title, err := sut.pageTitle(context.Background(), []byte(content))

	assert.NoError(t, err)
	assert.Empty(t, title)
}

func Test_page_title_returns_error_if_title_node_not_found(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">
	<head>
	</head>
	<body>
		
	</body>
	</html>`
	title, err := sut.pageTitle(context.Background(), []byte(content))

	assert.Empty(t, title)
	assert.ErrorContains(t, err, "title node not found")
}

func Test_page_title_returns_error_if_title_element_is_not_found_in_head_tag(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">
	<head>
	</head>
	<body>
		<title>Invalid title in Body</title>
	</body>
	</html>`
	title, err := sut.pageTitle(context.Background(), []byte(content))

	assert.Empty(t, title)
	assert.ErrorContains(t, err, "title node not found")
}

func Test_page_title_returns_error_if_head_element_is_missing_in_the_document(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">
	<body>
		<title>Invalid title in Body</title>
	</body>
	</html>`
	title, err := sut.pageTitle(context.Background(), []byte(content))

	assert.Empty(t, title)
	assert.ErrorContains(t, err, "head element not found in the document")
}
