package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/typelate/dom/domtest"
)

// TestGreetingPage asserts on a whole document with no test framework.
// Every result is checked because a domtest helper returns nil after reporting a failure.
func TestGreetingPage(t *testing.T) {
	var storage Storage
	mux := http.NewServeMux()
	routes(mux, &storage)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", res.StatusCode, http.StatusOK)
	}

	doc := domtest.ParseResponseDocument(t, res)
	if doc == nil {
		t.Fatal("could not parse the response")
	}

	h1 := doc.QuerySelector("h1.default-value")
	if h1 == nil {
		t.Fatal("expected the greeting to be marked as the default value")
	}
	if got, want := h1.TextContent(), "Hello, world!"; got != want {
		t.Errorf("h1 = %q, want %q", got, want)
	}

	input := doc.QuerySelector(`form input[name="name"]`)
	if input == nil {
		t.Fatal("expected a name field inside the form")
	}
	if got := input.GetAttribute("value"); got != "" {
		t.Errorf("value = %q, want it empty", got)
	}
}
