# example

A greeting server that one form field changes.
It exists so the [domtest README](../README.md) can point at code that runs, and so the tests beside it have something real to assert on.

```
go run ./domtest/example   # serves on :8080
go test ./domtest/example
```

## The server

[main.go](./main.go) registers the routes and [templates.gohtml](./templates.gohtml) renders them.

| route           | answers with                                                                   |
|-----------------|--------------------------------------------------------------------------------|
| `GET /`         | the whole page                                                                 |
| `POST /`        | the whole page, after storing the submitted name                               |
| `GET /greeting` | the heading alone, the kind of fragment a hypermedia library swaps into a page |

An empty name greets the world and marks the heading `default-value`.
Any other name is greeted by name and marks it `custom-value`.
That class is the useful part, because it lets a test select the state directly instead of matching the greeting text, so the assertion survives a copy edit.

`routes` takes a `StoreLoader` rather than the concrete `Storage`, which is the seam a fake slots into.

## The tests

[main_test.go](./main_test.go) asserts on the result with the standard library alone.
[table_test.go](./table_test.go) does the same work table-driven with testify.

| To see                                                      | Read                                                        |
|-------------------------------------------------------------|-------------------------------------------------------------|
| A whole document parsed with no test framework              | `TestGreetingPage` in [main_test.go](./main_test.go)        |
| Given, When, and Then cases with testify                    | `TestRoutes` in [table_test.go](./table_test.go)            |
| A form's method, action, and field values asserted          | the `the form round-trips the stored name` case             |
| A marker class selected instead of the page's wording       | the `a stored name is greeted instead` case                 |
| A fragment parsed against its parent element                | the `the greeting route answers with a fragment` case       |
| What the handler stored asserted alongside what it rendered | the `submitting a name stores it and greets it` case        |

## The shape of a handler suite

`TestRoutes` names its three phases so each case reads as a sentence about the server.
`Given` arranges the world, `When` builds the request, `Then` asserts on the response.
The phase types and `runCase` live inside the test function rather than in package scope, so another suite in the same package can reuse the names.
`runCase` owns the setup for every case, which leaves the table readable as a list of behaviors.
Each case gets a fresh mux over fresh storage, both local to the call, so cases stay independent and there is nothing to tear down.

Where a value lives follows from how many cases need it.
Anything every case supplies or returns belongs in the `Case` field signatures, which is why `When` returns a `*http.Request` and `Then` receives a `*http.Response`.
Anything only some cases need belongs in the `Given`, `When`, and `Then` structs.
A case wanting one more value gets a field, not a new parameter on every case that does not need it.
Naming each `Case` field after the struct it receives keeps the coupling easier to see.

Handing the storage to `Then` as well as `Given` lets a case assert on what the handler wrote.
Both phases hold an interface, so a suite with a real backend swaps in a fake.
In practice, you usually want to use a static type as the field type for phase fields.
Don't introduce interface types just for tests.
