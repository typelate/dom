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

func TestCase(t *testing.T) {
	type (
		Case struct {
			Name  string
			Given func(*testing.T, *fake.App)
			When  func(*testing.T) *http.Request
			Then  func(*testing.T, *http.Response, *fake.App)
		}
	)

	for _, tt := range []Case{
		{
			Name: "viewing the home page",
			Given: func(t *testing.T, app *fake.App) {
				app.ArticleReturns(blog.Article{
					Title:   "Greetings!",
					Content: "Hello, friends!",
					Error:   nil,
				})
			},
			When: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/1", nil)
			},
			Then: func(t *testing.T, response *http.Response, app *fake.App) {
				document := domtest.ParseResponseDocument(t, response)
				require.Equal(t, 1, app.ArticleArgsForCall(0))
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
			Given: func(t *testing.T, app *fake.App) { // GivenPtr removes some of the boilerplate in the block
				app.ArticleReturns(blog.Article{
					Error: fmt.Errorf("lemon"),
				})
			},
			When: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/1", nil)
			},
			Then: func(t *testing.T, response *http.Response, app *fake.App) {
				document := domtest.ParseResponseDocument(t, response)
				if msg := document.QuerySelector("#error-message"); assert.NotNil(t, msg) {
					require.Equal(t, "lemon", msg.TextContent())
				}
			},
		},
		{
			Name: "the page has an error and is requested by HTMX",
			Given: func(t *testing.T, app *fake.App) {
				app.ArticleReturns(blog.Article{
					Error: fmt.Errorf("lemon"),
				})
			},
			When: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/article/1", nil)
				req.Header.Set("HX-Request", "true")
				return req
			},
			Then: func(t *testing.T, response *http.Response, app *fake.App) {
				fragment := domtest.ParseResponseDocumentFragment(t, response, atom.Div)
				el := fragment.FirstElementChild()
				require.Equal(t, "lemon", el.TextContent())
				require.Equal(t, "*errors.errorString", el.GetAttribute("data-type"))
			},
		},
		{
			Name: "when the id is not an integer",
			When: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/banana", nil)
			},
			Then: func(t *testing.T, res *http.Response, f *fake.App) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			app := new(fake.App)

			if tt.Given != nil {
				tt.Given(t, app)
			}

			req := tt.When(t)

			mux := http.NewServeMux()
			blog.TemplateRoutes(mux, app)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			tt.Then(t, rec.Result(), app)
		})
	}
}
