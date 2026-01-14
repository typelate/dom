package dom

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom/spec"
)

var _ spec.Element = (*Element)(nil)

func TestElement_NodeType(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	document := Element{
		node: parsedDocument.FirstChild.NextSibling,
	}

	assert.Equal(t, spec.NodeTypeElement, document.NodeType())
}

func TestElement_IsConnected(t *testing.T) {
	t.Run("connected", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body>Hello!</body>
</html>`
		_, body := parseDocument(t, textHTML, "body")

		require.True(t, body.IsConnected())
	})

	t.Run("not connected", func(t *testing.T) {
		el := &Element{
			node: &html.Node{},
		}
		require.False(t, el.IsConnected())
	})
}

func TestElement_OwnerDocument(t *testing.T) {
	t.Run("connected", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body>Hello!</body>
</html>`
		document, body := parseDocument(t, textHTML, "body")

		owner := body.OwnerDocument()
		require.NotNil(t, owner)

		require.True(t, owner.(*Document).node == document.node)
	})
	t.Run("not connected", func(t *testing.T) {
		textNode := &Element{
			node: &html.Node{},
		}
		require.Nil(t, textNode.OwnerDocument())
	})
}

func TestElement_Length(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body>
	<div id='target'><br id='1'>2<div id='3'><span id='does-not-count'></span></div></div>
</body>
</html>`
	_, body := parseDocument(t, textHTML, "#target")
	assert.Equal(t, 3, body.Length())
}

func TestElement_ParentNode(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body><div id='app'></div></body>
</html>`
	_, app := parseDocument(t, textHTML, "html")

	parent := app.ParentNode()
	require.NotNil(t, parent)
	require.Equal(t, spec.NodeTypeDocument, parent.NodeType())
	_, ok := parent.(spec.Document)
	require.True(t, ok)
}

func TestElement_ParentElement(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body><div id='app'></div></body>
</html>`
	_, app := parseDocument(t, textHTML, "#app")

	parent, ok := app.ParentElement().(*Element)
	require.True(t, ok)
	assert.Equal(t, "body", parent.node.Data)
}

func TestElement_PreviousSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body><span id='oldest'></span><div id='target'></div></body>
</html>`
	_, target := parseDocument(t, textHTML, "#target")
	sibling := target.PreviousSibling()
	siblingNode, ok := sibling.(*Element)
	require.Truef(t, ok, "wrong type %T", sibling)

	require.Equal(t, "oldest", getAttribute(siblingNode.node, "id"))
}

func TestElement_NextSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body><div id='target'></div><span id='youngest'></span></body>
</html>`
	_, target := parseDocument(t, textHTML, "#target")
	sibling := target.NextSibling()
	siblingNode, ok := sibling.(*Element)
	require.Truef(t, ok, "wrong type %T", sibling)

	require.Equal(t, "youngest", getAttribute(siblingNode.node, "id"))
}

func TestElement_TextContent(t *testing.T) {
	t.Run("body no text", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Greetings!</title></head>
<body><div id='target'><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "body")
		text := target.TextContent()
		require.Equal(t, "\n", text)
	})

	t.Run("element no text", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Greetings!</title></head>
<body><div id='target'><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		text := target.TextContent()
		require.Equal(t, "", text)
	})

	t.Run("multiple lines with whitespace", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Greetings!</title></head>
<body>
<div id='target'>a
b<span>c </span>d
	e
</div>
</body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		text := target.TextContent()
		require.Equal(t, `a
bc d
	e
`, text)
	})
}

func TestElement_CloneNode(t *testing.T) {
	t.Run("not deep", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body><div id='target'><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		clonedElement, ok := target.CloneNode(false).(*Element)
		require.True(t, ok)
		require.Nil(t, clonedElement.node.FirstChild)
		require.Nil(t, clonedElement.node.Parent)
	})

	t.Run("deep", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title></title></head>
<body><div id='target'><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		clonedElement, ok := target.CloneNode(true).(*Element)
		require.True(t, ok)
		require.NotNil(t, clonedElement.node.FirstChild)
		require.Nil(t, clonedElement.node.Parent)
		require.EqualValues(t, clonedElement.node, clonedElement.node.FirstChild.Parent)
	})
}

func TestElement_IsSameNode(t *testing.T) {
	t.Run("different nodes", func(t *testing.T) {
		e1 := &Element{node: &html.Node{Type: html.ElementNode}}
		e2 := &Element{node: &html.Node{Type: html.ElementNode}}
		assert.False(t, e1.IsSameNode(e2))
	})
	t.Run("receiver node is nil", func(t *testing.T) {
		e1 := &Element{node: nil}
		e2 := &Element{node: &html.Node{Type: html.ElementNode}}
		assert.False(t, e1.IsSameNode(e2))
	})
	t.Run("param node is nil", func(t *testing.T) {
		e1 := &Element{node: &html.Node{Type: html.ElementNode}}
		e2 := &Element{node: nil}
		assert.False(t, e1.IsSameNode(e2))
	})
	t.Run("param is nil", func(t *testing.T) {
		e1 := &Element{node: &html.Node{Type: html.ElementNode}}
		assert.False(t, e1.IsSameNode(nil))
	})
	t.Run("true", func(t *testing.T) {
		node := &html.Node{Type: html.ElementNode}
		e1 := &Element{node: node}
		e2 := &Element{node: node}
		assert.True(t, e1.IsSameNode(e2))
	})
}

func TestElement_QuerySelector(t *testing.T) {
	parseFirstElement := func(t *testing.T, fragment string) *Element {
		nodes, err := html.ParseFragment(strings.NewReader(strings.TrimSpace(fragment)), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		node := NewNode(nodes[0])
		require.NotNil(t, node)
		el, ok := node.(*Element)
		require.True(t, ok)
		return el
	}
	t.Run("tree of nested results", func(t *testing.T) {
		element := parseFirstElement(t,
			/* language=html */ `<div id='n0'>
	<div id='n1'>
		<div id='n2'></div>
	</div>
	<div id='n3'></div>
</div>`)
		results := element.QuerySelectorAll("div")
		require.NotNil(t, results)
		assert.Equal(t, results.Length(), 3)

		for i := 0; i < results.Length(); i++ {
			result := results.Item(i)
			require.NotNil(t, result)
			assert.Equal(t, "DIV", result.TagName())
			assert.Equal(t, "n"+strconv.Itoa(i+1), result.ID())
		}
	})
}

func TestElement_CompareDocumentPosition(t *testing.T) {
	t.Run("same node", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><body><div id='x'></div></body></html>`
		_, x := parseDocument(t, textHTML, "#x")
		assert.Equal(t, spec.DocumentPosition(0), x.CompareDocumentPosition(x))
	})

	t.Run("a contains b", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><body><div id='a'><span id='b'></span></div></body></html>`
		doc, a := parseDocument(t, textHTML, "#a")
		b := doc.QuerySelector("#b")
		pos := a.CompareDocumentPosition(b)
		assert.Equal(t, spec.DocumentPositionContainedBy|spec.DocumentPositionFollowing, pos)
	})

	t.Run("b contains a", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><body><div id='a'><span id='b'></span></div></body></html>`
		doc, a := parseDocument(t, textHTML, "#a")
		b := doc.QuerySelector("#b")
		pos := b.CompareDocumentPosition(a)
		assert.Equal(t, spec.DocumentPositionContains|spec.DocumentPositionPreceding, pos)
	})

	t.Run("preceding", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><body><div id='a'></div><div id='b'></div></body></html>`
		doc, a := parseDocument(t, textHTML, "#a")
		b := doc.QuerySelector("#b")
		pos := a.CompareDocumentPosition(b)
		assert.Equal(t, spec.DocumentPositionFollowing, pos)
	})

	t.Run("following", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><body><div id='a'></div><div id='b'></div></body></html>`
		doc, a := parseDocument(t, textHTML, "#a")
		b := doc.QuerySelector("#b")
		pos := b.CompareDocumentPosition(a)
		assert.Equal(t, spec.DocumentPositionPreceding, pos)
	})

	t.Run("disconnected", func(t *testing.T) {
		a := &Element{node: &html.Node{Type: html.ElementNode, Data: "div"}}
		b := &Element{node: &html.Node{Type: html.ElementNode, Data: "div"}}
		pos := a.CompareDocumentPosition(b)
		assert.Equal(t, spec.DocumentPositionDisconnected|spec.DocumentPositionImplementationSpecific, pos)
	})
}

func TestElement_SetAttribute(t *testing.T) {
	t.Run("set attribute only adds element if required", func(t *testing.T) {
		n := &html.Node{Type: html.ElementNode, Data: "div"}
		div := &Element{node: n}

		div.SetAttribute("data-key", "0")
		div.SetAttribute("data-key", "1")

		require.Len(t, n.Attr, 1)
	})
}
