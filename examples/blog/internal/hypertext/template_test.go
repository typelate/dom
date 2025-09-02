package hypertext_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom/domtest"
	"github.com/typelate/dom/examples/blog/internal/fake"
	"github.com/typelate/dom/examples/blog/internal/hypertext"
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
		hypertext.TemplateRoutes(mux, app)

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
				given.app.ArticleReturns(hypertext.Article{
					Title:   "Greetings!",
					Content: "Hello, friends!",
					Error:   nil,
				})
			},
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, hypertext.TemplateRoutePaths{}.Article(1), nil)
			},
			Then: func(t *testing.T, response *http.Response, then Then) {
				document := domtest.ParseResponseDocument(t, response)
				assert.Equal(t, 1, then.app.ArticleArgsForCall(0))
				if heading := document.QuerySelector("h1"); assert.NotNil(t, heading) {
					assert.Equal(t, "Greetings!", heading.TextContent())
				}
				if content := document.QuerySelector("p"); assert.NotNil(t, content) {
					assert.Equal(t, "Hello, friends!", content.TextContent())
				}
			},
		},
		{
			Name: "the page has an error",
			Given: func(t *testing.T, given Given) { // GivenPtr removes some of the boilerplate in the block
				given.app.ArticleReturns(hypertext.Article{
					Error: fmt.Errorf("lemon"),
				})
			},
			When: func(t *testing.T, when When) *http.Request {
				return httptest.NewRequest(http.MethodGet, hypertext.TemplateRoutePaths{}.Article(1), nil)
			},
			Then: func(t *testing.T, response *http.Response, then Then) {
				document := domtest.ParseResponseDocument(t, response)
				if msg := document.QuerySelector("#error-message"); assert.NotNil(t, msg) {
					assert.Equal(t, "lemon", msg.TextContent())
				}
			},
		},
		{
			Name: "the page has an error and is requested by HTMX",
			Given: func(t *testing.T, given Given) {
				given.app.ArticleReturns(hypertext.Article{
					Error: fmt.Errorf("lemon"),
				})
			},
			When: func(t *testing.T, _ When) *http.Request {
				req := httptest.NewRequest(http.MethodGet, hypertext.TemplateRoutePaths{}.Article(1), nil)
				req.Header.Set("HX-Request", "true")
				return req
			},
			Then: func(t *testing.T, response *http.Response, _ Then) {
				fragment := domtest.ParseResponseDocumentFragment(t, response, atom.Div)
				if el := fragment.FirstElementChild(); assert.NotNil(t, el) {
					assert.Equal(t, "lemon", el.TextContent())
					assert.Equal(t, "*errors.errorString", el.GetAttribute("data-type"))
				}
			},
		},
		{
			Name: "when the id is not an integer",
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/banana", nil)
			},
			Then: func(t *testing.T, res *http.Response, _ Then) {
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			run(t, tt)
		})
	}
}
