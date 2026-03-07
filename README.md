# dom [![Go Reference](https://pkg.go.dev/badge/github.com/typelate/dom.svg)](https://pkg.go.dev/github.com/typelate/dom)

Pure Go implementation of the [WHATWG DOM](https://dom.spec.whatwg.org) backed by [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html). CSS selectors are provided by [andybalholm/cascadia](https://github.com/andybalholm/cascadia).

## Packages

| Package | Description |
|---------|-------------|
| `dom` | Implements `Document`, `Element`, `Text`, and `DocumentFragment` using `html.Node`. |
| `spec` | Interfaces matching the WHATWG DOM spec. Shared by `dom` and `browser`. |
| `domtest` | Test helpers that parse HTML strings or `http.Response` bodies into `spec` types. |
| `browser` | **Experimental.** Implements `spec` interfaces via `syscall/js` for WASM. |

## Example

```go
func TestGreeting(t *testing.T) {
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, httptest.NewRequest("GET", "/", nil))

	doc := domtest.ParseResponseDocument(t, res.Result())
	heading := doc.QuerySelector("h1")
	if heading == nil {
		t.Fatal("expected an h1")
	}
	if got := heading.TextContent(); got != "Hello" {
		t.Errorf("heading = %q, want %q", got, "Hello")
	}
}
```
