package dom_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/typelate/dom"
	"github.com/typelate/dom/spec"
)

func parseDoc(t *testing.T, in string) spec.Document {
	t.Helper()
	node, err := html.Parse(strings.NewReader(in))
	require.NoError(t, err)
	return dom.NewNode(node).(spec.Document)
}

// A deep clone must copy the fragment's children, not pad the result with nil
// entries. The nils used to survive into every method that walks d.nodes.
func TestDocumentFragment_CloneNode_deep_does_not_pad_with_nil(t *testing.T) {
	fragment := parseDocumentFragment(t, `<p>a</p><p>b</p>`)

	clone := fragment.CloneNode(true).(*dom.DocumentFragment)

	assert.Equal(t, 2, clone.ChildElementCount())
	assert.Equal(t, `<p>a</p><p>b</p>`, clone.String())
	assert.Equal(t, "ab", clone.TextContent())
	require.NotNil(t, clone.FirstElementChild())
	assert.Equal(t, "P", clone.FirstElementChild().TagName())
}

// GetElementsByTagName matches element names. A text node or comment whose data
// happens to read "div" is not a <div>.
func TestDocument_GetElementsByTagName_only_matches_elements(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><p>div</p><!--div--><div id="real"></div></body></html>`)

	list := document.GetElementsByTagName("div")

	require.Equal(t, 1, list.Length())
	assert.Equal(t, "real", list.Item(0).ID())
	assert.Equal(t, spec.NodeTypeElement, list.Item(0).NodeType())
}

func TestDocument_GetElementsByClassName_only_matches_elements(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><p class="c">hi</p><span>c</span></body></html>`)

	list := document.GetElementsByClassName("c")

	require.Equal(t, 1, list.Length())
	assert.Equal(t, "P", list.Item(0).TagName())
}

// An empty class name matches no classes, so it selects nothing.
func TestDocument_GetElementsByClassName_empty_name_matches_nothing(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><p class="x">hi</p></body></html>`)

	assert.Equal(t, 0, document.GetElementsByClassName("").Length())
	assert.Equal(t, 0, document.GetElementsByClassName("   ").Length())
}

// ReplaceChildren detaches every previous child, not just the first and last.
func TestElement_ReplaceChildren_detaches_every_previous_child(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><div id="p"><a>1</a><b>2</b><i>3</i></div></body></html>`)
	parent := document.QuerySelector("#p")
	middle := document.QuerySelector("b")
	require.True(t, middle.IsConnected())

	parent.ReplaceChildren()

	assert.Equal(t, "", parent.InnerHTML())
	assert.False(t, middle.IsConnected(), "detached child should not report as connected")
	assert.Nil(t, middle.ParentElement(), "detached child should have no parent element")
	assert.Nil(t, middle.PreviousSibling(), "detached child should have no siblings")
	assert.Nil(t, middle.NextSibling(), "detached child should have no siblings")
}

// Append and Prepend unwrap a DocumentFragment into its children; so must
// ReplaceChildren.
func TestElement_ReplaceChildren_unwraps_document_fragment(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><div id="p"><a>old</a></div></body></html>`)
	parent := document.QuerySelector("#p")
	fragment := parseDocumentFragment(t, `<span>new</span><em>also new</em>`)

	parent.ReplaceChildren(fragment)

	assert.Equal(t, `<span>new</span><em>also new</em>`, parent.InnerHTML())
}

// https://dom.spec.whatwg.org/#concept-node-length counts UTF-16 code units for
// character data, not bytes.
func TestText_Length_counts_utf16_code_units(t *testing.T) {
	for _, tt := range []struct {
		Data string
		Want int
	}{
		{Data: "hello", Want: 5},
		{Data: "héllo…", Want: 6},
		{Data: "a😀", Want: 3},
		{Data: "", Want: 0},
	} {
		t.Run(tt.Data, func(t *testing.T) {
			document := parseDoc(t, `<html><body><p>`+tt.Data+`</p></body></html>`)
			paragraph := document.QuerySelector("p")
			if tt.Data == "" {
				require.Nil(t, paragraph.FirstChild())
				return
			}
			text, ok := paragraph.FirstChild().(spec.Text)
			require.True(t, ok)

			assert.Equal(t, tt.Want, text.Length(), "Length(%q)", tt.Data)
		})
	}
}

func TestElementQueries_Contains_nil_is_false(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><p>x</p></body></html>`)

	assert.False(t, document.Contains(nil))
	assert.False(t, document.QuerySelector("body").Contains(nil))
}

// https://dom.spec.whatwg.org/#dom-document-createelement sets the node
// document of the new element to the document it was created from.
func TestDocument_created_nodes_have_an_owner_document(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body></body></html>`)

	element := document.CreateElement("div")
	require.NotNil(t, element.OwnerDocument())
	assert.True(t, element.OwnerDocument().IsSameNode(document))

	elementIs := document.CreateElementIs("div", "my-widget")
	require.NotNil(t, elementIs.OwnerDocument())
	assert.True(t, elementIs.OwnerDocument().IsSameNode(document))

	text := document.CreateTextNode("hi")
	require.NotNil(t, text.OwnerDocument())
	assert.True(t, text.OwnerDocument().IsSameNode(document))
}

// Once appended, the owner document comes from the tree, and it must still be
// the same document.
func TestDocument_created_element_keeps_owner_document_after_append(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body></body></html>`)
	element := document.CreateElement("div")

	document.Body().Append(element)

	assert.True(t, element.OwnerDocument().IsSameNode(document))
}

// Selectors are compiled with cascadia.MustCompile, so an invalid selector
// panics rather than returning an error.
func TestQuerySelector_panics_on_invalid_selector(t *testing.T) {
	// language=html
	document := parseDoc(t, `<html><body><p>x</p></body></html>`)

	assert.Panics(t, func() { document.QuerySelector("!!!") })
	assert.Panics(t, func() { document.QuerySelectorAll("!!!") })
	assert.Panics(t, func() { document.Body().Matches("!!!") })
	assert.Panics(t, func() { document.Body().Closest("!!!") })
}
