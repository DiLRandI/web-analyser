package webpage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/DiLRandI/web-analyser/internal/service/webpage/model"
	mc "github.com/DiLRandI/web-analyser/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func Test_login_form(t *testing.T) {
	testCases := []struct {
		desc   string
		expErr error
		expRes bool
		input  string
	}{
		{
			desc:   "hasLoginForm should return false if no login form found",
			expErr: nil,
			expRes: false,
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>no login form</title>
			</head>
			
			<body>
			<form>
			</body>
			
			</html>`,
		},
		{
			desc: `hasLoginForm should return true if login from found with <input type="submit"` +
				` value="Login">`,
			expErr: nil,
			expRes: true,
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>login form</title>
			</head>
			
			<body>
				<form action="test">
					<input type="text">
					<input type="password">
					<input type="submit" value="Login">
				</form>
			</body>
			
			</html>`,
		},
		{
			desc: `hasLoginForm should return true if login from found with ` +
				` <button type="submit">Login</button>`,
			expErr: nil,
			expRes: true,
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>login form</title>
			</head>
			
			<body>
				<form action="test">
					<input type="text">
					<input type="password">
					<button type="submit">Login</button>
				</form>
			</body>
			
			</html>`,
		},
		{
			desc:   `hasLoginForm should return false if the form is not a login form`,
			expErr: nil,
			expRes: false,
			input: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<title>Document</title>
			</head>
			<body>
				<form action="/singup">
					<button type="reset">X</button>
					<label for="Name">Name</label>
					<input type="text">
					<label for="Password">Password</label>
					<input type="password" name="password" id="password">
					<label for="Confirm password">Confirm password</label>
					<input type="password" name="confirmPassword" id="confirmPassword">
					<input type="submit" value="SignUp">
				</form>
			</body>
			</html>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			sut := &analyser{}
			actRes, actErr := sut.hasLoginForm(context.Background(), []byte(tc.input))

			if tc.expErr == nil {
				assert.NoError(t, actErr)
			} else {
				assert.ErrorIs(t, actErr, tc.expErr)
			}

			assert.Equal(t, tc.expRes, actRes)

		})
	}
}

func Test_link_details(t *testing.T) {
	testCases := []struct {
		desc       string
		hostUrl    string
		input      string
		expErr     error
		expRes     []*model.Link
		mockClient []struct {
			getUrl string
			res    *http.Response
			err    error
		}
	}{
		{
			desc:    "Should return empty array if no links found",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
			
			</html>`,
			expErr: nil,
			expRes: make([]*model.Link, 0),
		},
		{
			desc:    "Should return valid link object with internal true for fragment link",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="#id-test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "#id-test",
					IsInternal:     true,
					LinkStatus:     model.LinkStatusActive,
					HttpStatusCode: http.StatusOK,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: "http://www.test.com/#id-test",
					res: &http.Response{
						StatusCode: http.StatusOK,
					},
					err: nil,
				},
			},
		},
		{
			desc:    "Should return valid link object with internal true for path link",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="/test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "/test",
					IsInternal:     true,
					LinkStatus:     model.LinkStatusActive,
					HttpStatusCode: http.StatusOK,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: mock.Anything,
					res: &http.Response{
						StatusCode: http.StatusOK,
					},
					err: nil,
				},
			},
		},
		{
			desc:    "Should return valid link object with internal true for full path link",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="http://www.test.com/test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "http://www.test.com/test",
					IsInternal:     true,
					LinkStatus:     model.LinkStatusActive,
					HttpStatusCode: http.StatusOK,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: mock.Anything,
					res: &http.Response{
						StatusCode: http.StatusOK,
					},
					err: nil,
				},
			},
		},
		{
			desc:    "Should return valid link object with internal false for external link",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="http://www.different-test.com/test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "http://www.different-test.com/test",
					IsInternal:     false,
					LinkStatus:     model.LinkStatusActive,
					HttpStatusCode: http.StatusOK,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: mock.Anything,
					res: &http.Response{
						StatusCode: http.StatusOK,
					},
					err: nil,
				},
			},
		},
		{
			desc:    "Should return valid link object with internal false for external link",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="http://www.different-test.com/test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "http://www.different-test.com/test",
					IsInternal:     false,
					LinkStatus:     model.LinkStatusActive,
					HttpStatusCode: http.StatusOK,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: mock.Anything,
					res: &http.Response{
						StatusCode: http.StatusOK,
					},
					err: nil,
				},
			},
		},
		{
			desc:    "Should return link status inactive for not found url with correct status code",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="http://www.different-test.com/test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "http://www.different-test.com/test",
					IsInternal:     false,
					LinkStatus:     model.LinkStatusInactive,
					HttpStatusCode: http.StatusNotFound,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: mock.Anything,
					res: &http.Response{
						StatusCode: http.StatusNotFound,
					},
					err: nil,
				},
			},
		},
		{
			desc:    "Should return link status inactive with -1 status if request failed",
			hostUrl: "http://www.test.com",
			input: `<!DOCTYPE html>
			<html lang="en">
			
			<head>
				<title>Document</title>
			</head>
			
			<body>
			</body>
				<a href="http://www.different-test.com/test">Go to test</a>
			</html>`,
			expErr: nil,
			expRes: []*model.Link{
				{
					Name:           "Go to test",
					Url:            "http://www.different-test.com/test",
					IsInternal:     false,
					LinkStatus:     model.LinkStatusInactive,
					HttpStatusCode: -1,
				},
			},
			mockClient: []struct {
				getUrl string
				res    *http.Response
				err    error
			}{
				{
					getUrl: mock.Anything,
					res:    nil,
					err:    errors.New("Test failure"),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mc := new(mc.WebClientMock)
			for _, v := range tc.mockClient {
				mc.On("Get", v.getUrl).Return(v.res, v.err)
			}

			sut := &analyser{
				client: mc,
			}

			actRes, actErr := sut.linksDetail(context.Background(), tc.hostUrl, []byte(tc.input))

			if tc.expErr == nil {
				assert.NoError(t, actErr)
			} else {
				assert.ErrorIs(t, actErr, tc.expErr)
			}

			assert.Equal(t, tc.expRes, actRes)
		})
	}
}
