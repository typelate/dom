//go:build js

package browser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/typelate/dom/browser"
	"github.com/typelate/dom/spec"
)

// Append adds to the end and Prepend adds to the front. These used to be wired
// to each other's DOM call.
func TestElement_Append_and_Prepend(t *testing.T) {
	document := browser.OpenDocument()
	parent := document.CreateElement("div")

	first := document.CreateElement("span")
	first.SetAttribute("id", "first")
	last := document.CreateElement("span")
	last.SetAttribute("id", "last")

	parent.Append(last)
	parent.Prepend(first)

	children := parent.Children()
	require.Equal(t, 2, children.Length())
	assert.Equal(t, "first", children.Item(0).ID())
	assert.Equal(t, "last", children.Item(1).ID())
}

// The "is" option must be passed as a JavaScript object. Passing a Go struct
// made every call panic with "ValueOf: invalid value".
func TestDocument_CreateElementIs(t *testing.T) {
	document := browser.OpenDocument()

	var element spec.Element
	require.NotPanics(t, func() {
		element = document.CreateElementIs("div", "my-widget")
	})

	require.NotNil(t, element)
	assert.Equal(t, "DIV", element.TagName())
}

// https://dom.spec.whatwg.org/#concept-node-length is the number of children.
func TestElement_Length(t *testing.T) {
	document := browser.OpenDocument()
	parent := document.CreateElement("div")

	assert.Equal(t, 0, parent.(spec.ChildNode).Length())

	parent.AppendChild(document.CreateElement("span"))
	parent.AppendChild(document.CreateTextNode("text"))

	assert.Equal(t, 2, parent.(spec.ChildNode).Length())
}
