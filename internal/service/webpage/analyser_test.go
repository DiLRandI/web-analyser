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

func Test_page_heading_count_should_return_valid_heading_count(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">
	<body>
	</body>
	</html>`
	expected := map[string]int{
		"h1": 0,
		"h2": 0,
		"h3": 0,
		"h4": 0,
		"h5": 0,
		"h6": 0,
	}
	headings, err := sut.headingDetails(context.Background(), []byte(content))

	assert.NoError(t, err)
	assert.EqualValues(t, headings, expected)
}

func Test_page_heading_count_should_return_valid_heading_count_for_each_heading(t *testing.T) {
	sut := &analyser{}
	content := `
	<html lang="en">

<body>
    <h1>This is a h1</h1>
    <h2>This is a h2</h2>
    <h3>This is a h3</h3>
    <h4>This is a h4</h4>
    <h5>This is a h5</h5>
    <h6>This is a h6</h6>
</body>

</html>`
	expected := map[string]int{
		"h1": 1,
		"h2": 1,
		"h3": 1,
		"h4": 1,
		"h5": 1,
		"h6": 1,
	}
	headings, err := sut.headingDetails(context.Background(), []byte(content))

	assert.NoError(t, err)
	assert.EqualValues(t, headings, expected)
}
