# domtest [![Go Reference](https://pkg.go.dev/badge/github.com/typelate/dom/domtest.svg)](https://pkg.go.dev/github.com/typelate/dom/domtest)

domtest parses the HTML a handler wrote into `spec` types, so a test can query it with CSS selectors instead of matching substrings.
Nothing here needs a browser.

```go
import "github.com/typelate/dom/domtest"
```

A runnable example lives in [./example](./example), and its [README](./example/README.md) maps each idea below to the test that demonstrates it.

## Parsers

Six functions cover three input sources and two parse modes.

|              | `*http.Response`                | `string`                      | `io.Reader`                   |
|--------------|---------------------------------|-------------------------------|-------------------------------|
| **Document** | `ParseResponseDocument`         | `ParseStringDocument`         | `ParseReaderDocument`         |
| **Fragment** | `ParseResponseDocumentFragment` | `ParseStringDocumentFragment` | `ParseReaderDocumentFragment` |

Document parsers return `spec.Document`, fragment parsers return `spec.DocumentFragment`.
Both response parsers close `res.Body`, on the success path and on failure.

## What to assert on

The markup a hypermedia server returns is its interface. Forms, links, and fields carry what a user can do next, so assert on those.
A form that posts to a route the server serves, and a field that round-trips its value under the name the handler reads, are claims that survive a rewording of the page.

A `strings.Contains` check over the response body keeps passing after a field escapes its form or loses its `name`.
The selector `form input[name="name"]` stops matching the moment either happens. Name the container and the element, then assert on the attribute the next request depends on.

Marker classes work the same way.
The example puts `default-value` or `custom-value` on the heading depending on which branch produced it, so a test selects the state directly and survives a copy edit.

Attribute selectors accept hyphenated names, so `hx-*`, `fx-*`, and `data-*` are reachable like any other attribute.
Passing a control's target attribute back into `QuerySelector` checks in one line that what it points at exists.

Whole-page golden comparisons fail on every unrelated edit and say nothing about which affordance broke. Keep assertions narrow and let the selector carry the structural claim.

## Reusable assertions

A check worth making twice belongs in a helper, and `spec` already has the interface to write it against.
Take a `spec.ElementQueries` and the same helper runs over a whole document or any single element.

```go
// assertLinksResolve fails for every anchor whose href is missing, empty, or
// not something net/url will parse.
func assertLinksResolve(t *testing.T, root spec.ElementQueries) {
	t.Helper()
	for a := range root.QuerySelectorSequence("a") {
		href := a.GetAttribute("href")
		if !assert.NotEmptyf(t, href, "anchor %q has no href", a.TextContent()) {
			continue
		}
		_, err := url.Parse(href)
		assert.NoErrorf(t, err, "anchor %q has an unparsable href %q", a.TextContent(), href)
	}
}
```

Pass a document to sweep the page, or an element to narrow it, as in `assertLinksResolve(t, doc.QuerySelector("nav"))`.
A `spec.DocumentFragment` declares the query methods without the `Contains` and `GetElementsBy` group, so it does not satisfy `ElementQueries`.
Reach into a fragment first and hand the helper an element.

`url.Parse` rejects less than you would hope.
It turns down bad escapes like `%zz`, spaces in a host, and control characters, but accepts the empty string and `javascript:alert(1)`, which is why the emptiness check earns its keep.
A helper is where project rules go, such as requiring a leading slash on internal links.

## Fragments and the parent element

A handler answering a partial request returns a fragment, and it has to be parsed in the context of the element it will land in.
The parent argument names that element and decides the result.
HTML parsing rules discard a `<tr>` outside a table, so the wrong parent produces an empty fragment:

| parent       | result for `<tr><td>a</td></tr><tr><td>b</td></tr>` |
|--------------|-----------------------------------------------------|
| `atom.Body`  | 0 elements                                          |
| `atom.Tbody` | 2 elements                                          |

`html.ParseFragment` runs the same algorithm a browser runs for `innerHTML`, so a fragment that vanishes here would not have landed in the page either.
Pass the element the response gets swapped into and the test checks the swap contract for free.

Use `atom.Tbody` for table rows, `atom.Select` for options, and `atom.Body` or `atom.Div` for ordinary flow content.

## Failure behavior

A helper that cannot read or parse its input reports it with `t.Error` and returns nil. `Error` does not stop the test, so the assertions run on with a nil document and the first method call panics.
Check the result, or call `t.Failed`.
The nil is a true nil interface value, so comparing against `nil` is reliable.

Queries behave the same way, and a miss is nil rather than a zero value.
Checking for it is the one piece of ceremony this package asks for.

## With testify

Nil checks read better as `require`, which stops the test, and comparisons read better as `assert`, which lets one run report every mismatch.
Anything you then dereference needs `require`, because a nil `Element` panics on the next method call.
Anything only compared can be `assert`, so a broken page reports its heading, its class, and its fields together instead of one at a time.

The helpers follow the same convention, reporting with `t.Error` and returning nil, which pairs with `require.NotNil`.

## TestingT

The helpers accept a `TestingT` rather than `*testing.T`, which `*testing.T`, `*testing.B`, and `*testing.F` all satisfy.
The narrower interface also lets a test double record what the helpers reported, which is how this package tests itself.

## Notes

- A missing element is nil rather than a zero value. `QuerySelector` returns nil when nothing matches, while `QuerySelectorAll` returns a list of length 0.
- An invalid CSS selector panics, because queries compile with `cascadia.MustCompile`.
- `TextContent` joins descendant text with no separator and keeps the source whitespace, so an indented row of two cells reads as `"\n\tAda\n\t42\n"`. Wrap it in `strings.TrimSpace`.
- Document parsing inserts missing `html`, `head`, and `body` elements, the same as `html.Parse`, so `ParseStringDocument(t, "<p>hi</p>")` returns a document whose `Body().InnerHTML()` is `<p>hi</p>`.
- `spec.Document` leaves out the `ParentNode` child-management methods, so go through `Body()` when you need to walk children.
