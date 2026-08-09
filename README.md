# dom [![Go Reference](https://pkg.go.dev/badge/github.com/typelate/dom.svg)](https://pkg.go.dev/github.com/typelate/dom)

A pure Go implementation of a subset of the [WHATWG DOM](https://dom.spec.whatwg.org), backed by [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html), with CSS selectors from [andybalholm/cascadia](https://github.com/andybalholm/cascadia).

Its main use is asserting on the HTML your handlers return, without a browser.

```
go get github.com/typelate/dom
```

## Packages

| Package | Description |
|---------|-------------|
| `spec` | Interfaces for a subset of the WHATWG DOM. Implemented by `dom` and `browser`. |
| `dom` | `Document`, `Element`, `Text`, and `DocumentFragment` over `html.Node`. |
| `domtest` | Parses HTML strings, `io.Reader`s, and `http.Response` bodies into `spec` types. |
| `browser` | **Experimental.** The same interfaces over the real browser DOM via `syscall/js`, for WASM. |

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

To assert on a fragment rather than a whole page, parse it in the context of the
element it will land in — the parent decides how the markup is interpreted:

```go
rows := domtest.ParseResponseDocumentFragment(t, res.Result(), atom.Tbody)
```

## Notes

- An invalid CSS selector panics, because queries compile with `cascadia.MustCompile`.
- `domtest` helpers report failures with `t.Error` and return `nil`, which does not stop the test. Check the result before using it.
- Collections are not uniformly live: `Children` reflects later mutations, while `GetElementsByTagName` and `GetElementsByClassName` return a snapshot.
