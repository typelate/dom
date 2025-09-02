package blog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom/domtest"
	"github.com/typelate/dom/domtest/example/blog"
	"github.com/typelate/dom/domtest/example/blog/internal/fake"
)

func TestAPI(t *testing.T) {
	type (
		Given struct {
			app *fake.App
		}
		When struct{}
		Then struct {
			app *fake.App
		}
		Case struct {
			Name  string
			Given func(*testing.T, Given)
			When  func(*testing.T, When) *http.Request
			Then  func(*testing.T, *http.Response, Then)
		}
	)

	run := func(t *testing.T, tc Case) {
		app := new(fake.App)

		if tc.Given != nil {
			tc.Given(t, Given{
				app: app,
			})
		}

		req := tc.When(t, When{})

		mux := http.NewServeMux()
		blog.TemplateRoutes(mux, app)

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		tc.Then(t, rec.Result(), Then{
			app: app,
		})
	}

	for _, tt := range []Case{
		{
			Name: "viewing the home page",
			Given: func(t *testing.T, given Given) {
				given.app.ArticleReturns(blog.Article{
					Title:   "Greetings!",
					Content: "Hello, friends!",
					Error:   nil,
				})
			},
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/1", nil)
			},
			Then: func(t *testing.T, response *http.Response, then Then) {
				document := domtest.ParseResponseDocument(t, response)
				require.Equal(t, 1, then.app.ArticleArgsForCall(0))
				if heading := document.QuerySelector("h1"); assert.NotNil(t, heading) {
					require.Equal(t, "Greetings!", heading.TextContent())
				}
				if content := document.QuerySelector("p"); assert.NotNil(t, content) {
					require.Equal(t, "Hello, friends!", content.TextContent())
				}
			},
		},
		{
			Name: "the page has an error",
			Given: func(t *testing.T, given Given) { // GivenPtr removes some of the boilerplate in the block
				given.app.ArticleReturns(blog.Article{
					Error: fmt.Errorf("lemon"),
				})
			},
			When: func(t *testing.T, when When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/1", nil)
			},
			Then: func(t *testing.T, response *http.Response, then Then) {
				document := domtest.ParseResponseDocument(t, response)
				if msg := document.QuerySelector("#error-message"); assert.NotNil(t, msg) {
					require.Equal(t, "lemon", msg.TextContent())
				}
			},
		},
		{
			Name: "the page has an error and is requested by HTMX",
			Given: func(t *testing.T, given Given) {
				given.app.ArticleReturns(blog.Article{
					Error: fmt.Errorf("lemon"),
				})
			},
			When: func(t *testing.T, _ When) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/article/1", nil)
				req.Header.Set("HX-Request", "true")
				return req
			},
			Then: func(t *testing.T, response *http.Response, _ Then) {
				fragment := domtest.ParseResponseDocumentFragment(t, response, atom.Div)
				el := fragment.FirstElementChild()
				require.Equal(t, "lemon", el.TextContent())
				require.Equal(t, "*errors.errorString", el.GetAttribute("data-type"))
			},
		},
		{
			Name: "when the id is not an integer",
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/banana", nil)
			},
			Then: func(t *testing.T, res *http.Response, _ Then) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			run(t, tt)
		})
	}
}
