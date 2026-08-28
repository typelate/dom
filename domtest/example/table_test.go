package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom/domtest"
)

// TestRoutes drives every route through the same three phases.
func TestRoutes(t *testing.T) {
	type (
		Given struct{ storage StoreLoader }
		When  struct{}
		Then  struct{ storage StoreLoader }
		Case  struct {
			Name  string
			Given func(*testing.T, Given)
			When  func(*testing.T, When) *http.Request
			Then  func(*testing.T, Then, *http.Response)
		}
	)

	runCase := func(t *testing.T, tc Case) {
		var storage Storage
		if tc.Given != nil {
			tc.Given(t, Given{storage: &storage})
		}

		mux := http.NewServeMux()
		routes(mux, &storage)

		req := tc.When(t, When{})
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if tc.Then != nil {
			tc.Then(t, Then{storage: &storage}, rec.Result())
		}
	}

	for _, tc := range []Case{
		{
			Name: "an empty name greets the world",
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			Then: func(t *testing.T, _ Then, res *http.Response) {
				assert.Equal(t, http.StatusOK, res.StatusCode)

				doc := domtest.ParseResponseDocument(t, res)
				require.NotNil(t, doc)

				h1 := doc.QuerySelector("h1")
				require.NotNil(t, h1)
				assert.Equal(t, "Hello, world!", h1.TextContent())
				assert.Equal(t, "default-value", h1.ClassName())
			},
		},
		{
			Name: "a stored name is greeted instead",
			Given: func(t *testing.T, g Given) {
				g.storage.Store("friend")
			},
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			Then: func(t *testing.T, _ Then, res *http.Response) {
				assert.Equal(t, http.StatusOK, res.StatusCode)

				doc := domtest.ParseResponseDocument(t, res)
				require.NotNil(t, doc)

				h1 := doc.QuerySelector("h1.custom-value")
				require.NotNil(t, h1, "the greeting should be marked as a custom value")
				assert.Equal(t, "Hello, friend!", h1.TextContent())

				assert.Nil(t, doc.QuerySelector("h1.default-value"))
			},
		},
		{
			Name: "the form round-trips the stored name",
			Given: func(t *testing.T, g Given) {
				g.storage.Store("friend")
			},
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			Then: func(t *testing.T, _ Then, res *http.Response) {
				doc := domtest.ParseResponseDocument(t, res)
				require.NotNil(t, doc)

				form := doc.QuerySelector("form")
				require.NotNil(t, form)
				assert.Equal(t, "POST", form.GetAttribute("method"))
				assert.Equal(t, "/", form.GetAttribute("action"))

				input := doc.QuerySelector(`form input[name="name"]`)
				require.NotNil(t, input)
				assert.Equal(t, "friend", input.GetAttribute("value"))
			},
		},
		{
			Name: "submitting a name stores it and greets it",
			When: func(t *testing.T, _ When) *http.Request {
				body := url.Values{"name": []string{"friend"}}.Encode()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			Then: func(t *testing.T, then Then, res *http.Response) {
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "friend", then.storage.Load())

				doc := domtest.ParseResponseDocument(t, res)
				require.NotNil(t, doc)

				h1 := doc.QuerySelector("h1.custom-value")
				require.NotNil(t, h1)
				assert.Equal(t, "Hello, friend!", h1.TextContent())
			},
		},
		{
			Name: "submitting an empty name falls back to the default",
			Given: func(t *testing.T, g Given) {
				g.storage.Store("friend")
			},
			When: func(t *testing.T, _ When) *http.Request {
				body := url.Values{"name": []string{""}}.Encode()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			Then: func(t *testing.T, then Then, res *http.Response) {
				assert.Equal(t, "", then.storage.Load())

				doc := domtest.ParseResponseDocument(t, res)
				require.NotNil(t, doc)

				h1 := doc.QuerySelector("h1.default-value")
				require.NotNil(t, h1)
				assert.Equal(t, "Hello, world!", h1.TextContent())
			},
		},
		{
			Name: "the greeting route answers with a fragment",
			Given: func(t *testing.T, g Given) {
				g.storage.Store("friend")
			},
			When: func(t *testing.T, _ When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/greeting", nil)
			},
			Then: func(t *testing.T, _ Then, res *http.Response) {
				fragment := domtest.ParseResponseDocumentFragment(t, res, atom.Body)
				require.NotNil(t, fragment)

				require.Equal(t, 1, fragment.ChildElementCount())
				h1 := fragment.QuerySelector("h1")
				require.NotNil(t, h1)
				assert.Equal(t, "custom-value", h1.ClassName())
				assert.Equal(t, "Hello, friend!", h1.TextContent())
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc)
		})
	}
}
