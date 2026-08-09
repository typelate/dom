// Package domtest provides test helpers that parse HTML into spec types.
//
// Every helper reports parse and read failures with TestingT.Error and returns
// nil. Error does not stop the test, so check the result or call t.Failed
// before dereferencing it.
package domtest

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom"
	"github.com/typelate/dom/spec"
)

// TestingT is the subset of *testing.T the helpers use.
type TestingT interface {
	Helper()
	Error(...any)
	Log(...any)
	Errorf(format string, args ...any)
	FailNow()
	Failed() bool
	SkipNow()
}

// ParseResponseDocument parses a whole HTML document from res and closes its
// body.
func ParseResponseDocument(t TestingT, res *http.Response) spec.Document {
	t.Helper()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
		return nil
	}
	if err := res.Body.Close(); err != nil {
		t.Error(err)
		return nil
	}
	return ParseReaderDocument(t, bytes.NewReader(buf))
}

// ParseStringDocument parses a whole HTML document from s.
func ParseStringDocument(t TestingT, s string) spec.Document {
	t.Helper()
	return ParseReaderDocument(t, strings.NewReader(s))
}

// ParseReaderDocument parses a whole HTML document from r. Missing html, head,
// and body elements are inserted, as by html.Parse.
func ParseReaderDocument(t TestingT, r io.Reader) spec.Document {
	t.Helper()
	node, err := html.Parse(r)
	if err != nil {
		t.Error(err)
		return nil
	}
	return dom.NewNode(node).(spec.Document)
}

// ParseResponseDocumentFragment parses an HTML fragment from res and closes its
// body. See ParseReaderDocumentFragment for the meaning of parent.
func ParseResponseDocumentFragment(t TestingT, res *http.Response, parent atom.Atom) spec.DocumentFragment {
	t.Helper()
	defer closeAndCheckError(t, res.Body)
	return ParseReaderDocumentFragment(t, res.Body, parent)
}

// ParseStringDocumentFragment parses an HTML fragment from in. See
// ParseReaderDocumentFragment for the meaning of parent.
func ParseStringDocumentFragment(t TestingT, in string, parent atom.Atom) spec.DocumentFragment {
	t.Helper()
	return ParseReaderDocumentFragment(t, strings.NewReader(in), parent)
}

// ParseReaderDocumentFragment parses an HTML fragment from r as if it were
// inside a parent element, which decides how the markup is interpreted: use
// atom.Tbody for rows, atom.Select for options, and atom.Body or atom.Div for
// ordinary flow content.
func ParseReaderDocumentFragment(t TestingT, r io.Reader, parent atom.Atom) spec.DocumentFragment {
	t.Helper()

	body, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
		return nil
	}
	nodes, err := html.ParseFragment(bytes.NewReader(body), &html.Node{
		Type:     html.ElementNode,
		Data:     parent.String(),
		DataAtom: parent,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	return dom.NewDocumentFragment(nodes)
}

func closeAndCheckError(t TestingT, c io.Closer) {
	if err := c.Close(); err != nil {
		t.Error(err)
	}
}
