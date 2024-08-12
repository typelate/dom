package dom

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/spec"
)

func TestElement_NodeType(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
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
<html lang="us-en">
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
<html lang="us-en">
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
<html lang="us-en">
<head><title></title></head>
<body>
	<div id="target"><br id="1">2<div id="3"><span id="does-not-count"></span></div></div>
</body>
</html>`
	_, body := parseDocument(t, textHTML, "#target")
	assert.Equal(t, 3, body.Length())
}

func TestElement_ParentNode(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><div id="app"></div></body>
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
<html lang="us-en">
<head><title></title></head>
<body><div id="app"></div></body>
</html>`
	_, app := parseDocument(t, textHTML, "#app")

	parent, ok := app.ParentElement().(*Element)
	require.True(t, ok)
	assert.Equal(t, "body", parent.node.Data)
}

func TestElement_PreviousSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><span id="oldest"></span><div id="target"></div></body>
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
<html lang="us-en">
<head><title></title></head>
<body><div id="target"></div><span id="youngest"></span></body>
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
<html lang="us-en">
<head><title>Greetings!</title></head>
<body><div id="target"><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "body")
		text := target.TextContent()
		require.Equal(t, "\n", text)
	})

	t.Run("element no text", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title>Greetings!</title></head>
<body><div id="target"><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		text := target.TextContent()
		require.Equal(t, "", text)
	})

	t.Run("multiple lines with whitespace", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title>Greetings!</title></head>
<body>
<div id="target">a
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
<html lang="us-en">
<head><title></title></head>
<body><div id="target"><span></span></div></body>
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
<html lang="us-en">
<head><title></title></head>
<body><div id="target"><span></span></div></body>
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
			/* language=html */ `<div id="n0">
	<div id="n1">
		<div id="n2"></div>
	</div>
	<div id="n3"></div>
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

func Test_elementList_Forwards(t *testing.T) {
	//
	documentAST, err := html.Parse(strings.NewReader(
		/* language=html */
		`<!DOCTYPE html>
<html lang="en-us">
<head>
    <meta charset="utf-8">
    <title>Test Data</title>
</head>
<body>
    <a href="/1">1</a>
	<a href="/2">2</a>
	<div>
		<a href="/3">3</a>
	</div>
	<a href="/4">4</a>
</body>
</html>`))
	require.NoError(t, err)

	var document spec.Document = &Document{node: documentAST}

	callCount := 0
	for index, element := range document.GetElementsByTagName(`a`).Each {
		assert.Equal(t, "/"+strconv.Itoa(index+1), element.GetAttribute("href"))
		callCount++
	}
	assert.Equal(t, 4, callCount)
}

func Test_querySelector_Forwards(t *testing.T) {
	documentAST, err := html.Parse(strings.NewReader(
		/* language=html */
		`<!DOCTYPE html>
<html lang="en-us">
<head>
    <meta charset="utf-8">
    <title>Test Data</title>
</head>
<body>
    <div data-index="0">1</div>
	<div data-index="1" >2</div>
	<span>
		<div data-index="2" >2</div>
	</span>
	<div data-index="3" >3</div>
</body>
</html>`))
	require.NoError(t, err)

	var document spec.Document = &Document{node: documentAST}

	callCount := 0
	for index, element := range document.QuerySelectorAll(`[data-index]`).Each {
		assert.Equal(t, strconv.Itoa(index), element.GetAttribute("data-index"))
		callCount++
	}
	assert.Equal(t, 4, callCount)
}

func Test_siblingList_Forwards(t *testing.T) {
	documentAST, err := html.Parse(strings.NewReader(
		/* language=html */
		`<!DOCTYPE html>
<html lang="en-us">
<head>
    <meta charset="utf-8">
    <title>Test Data</title>
</head>
<body>
    <span data-index="0">1</span>
	<div data-index="1" >2</div>
	<a data-index="2" href="/three">3</a>
</body>
</html>`))
	require.NoError(t, err)

	var document spec.Document = &Document{node: documentAST}

	callCount := 0
	for index, element := range document.QuerySelector(`body`).Children().Each {
		assert.Equal(t, strconv.Itoa(index), element.GetAttribute("data-index"))
		callCount++
	}
	assert.Equal(t, 3, callCount)
}
