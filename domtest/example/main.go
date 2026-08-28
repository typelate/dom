// Command example serves a greeting that one form field changes. It exists so
// the domtest README can point at code that runs, and so the tests beside it
// have something real to assert on.
//
// The server has three routes. GET / and POST / answer with a whole document,
// while GET /greeting answers with the heading alone, the kind of fragment a
// hypermedia library swaps into a page. An empty name greets the world and
// marks the heading default-value. Any other name is greeted by name and marks
// the heading custom-value, so a test selects the state by class instead of by
// wording.
//
// main_test.go asserts on the result with the standard library alone.
// table_test.go does the same work table-driven with testify.
package main

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"sync"
)

//go:embed templates.gohtml
var templates embed.FS

func main() {
	var s Storage

	mux := http.NewServeMux()
	routes(mux, &s)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// greeting is what the templates render. Class marks whether the page fell
// back to the default name, which gives a test something to select on without
// matching the greeting text itself.
type greeting struct {
	Name    string
	Greeted string
	Class   string
}

func newGreeting(storage string) greeting {
	if storage == "" {
		return greeting{Greeted: "world", Class: "default-value"}
	}
	return greeting{Name: storage, Greeted: storage, Class: "custom-value"}
}

// StoreLoader holds the name the greeting uses. routes takes this interface
// rather than a concrete type, so a test can hand it another implementation.
type StoreLoader interface {
	Store(string)
	Load() string
}

// Storage is an in-memory StoreLoader that is safe for concurrent use. Its
// zero value is ready to use and loads the empty name.
type Storage struct {
	m     sync.Mutex
	value string
}

// Store replaces the stored name.
func (s *Storage) Store(val string) {
	s.m.Lock()
	s.value = val
	s.m.Unlock()
}

// Load returns the stored name, which is empty until Store is called.
func (s *Storage) Load() string {
	s.m.Lock()
	val := s.value
	s.m.Unlock()
	return val
}

// routes registers the handlers on mux. It takes a StoreLoader rather than a
// *Storage so a test can supply its own.
func routes(mux *http.ServeMux, storage StoreLoader) {
	pages := template.Must(template.ParseFS(templates, "templates.gohtml"))

	// render takes a callback so that each call to ExecuteTemplate
	// can map a template name to a static type. This is helpful if you are using `muxt check`.
	// So that you don't need to write type-checks for your template actions.
	render := func(w http.ResponseWriter, execute func(w io.Writer) error) {
		var buf bytes.Buffer
		if err := execute(&buf); err != nil {
			http.Error(w, "could not render the page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		render(w, func(w io.Writer) error {
			return pages.ExecuteTemplate(w, "page", newGreeting(storage.Load()))
		})
	})

	mux.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
		storage.Store(r.FormValue("name"))
		render(w, func(w io.Writer) error {
			return pages.ExecuteTemplate(w, "page", newGreeting(storage.Load()))
		})
	})

	mux.HandleFunc("GET /greeting", func(w http.ResponseWriter, _ *http.Request) {
		render(w, func(w io.Writer) error {
			return pages.ExecuteTemplate(w, "greeting", newGreeting(storage.Load()))
		})
	})
}
