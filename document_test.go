package dom

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/html"

	"github.com/typelate/dom/spec"
)

func TestDocument_NodeType(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
	document, _ := parseDocument(t, textHTML, "")

	assert.Equal(t, spec.NodeTypeDocument, document.NodeType())
}

func TestDocument_CloneNode(t *testing.T) {
	t.Run("deep", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
		document, _ := parseDocument(t, textHTML, "")

		result := document.CloneNode(true)

		copied, ok := result.(*Document)
		assert.True(t, ok, "correct type")
		assert.False(t, copied.node == document.node, "not the same *html.NewNode")
		assert.NotNil(t, copied.node.FirstChild)
		assert.False(t, copied.node.FirstChild == document.node.FirstChild, "not the same child *html.NewNode")
	})
	t.Run("not deep", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
		document, _ := parseDocument(t, textHTML, "")

		result := document.CloneNode(false)

		copied, ok := result.(*Document)
		assert.True(t, ok, "correct type")
		assert.Nil(t, copied.node.FirstChild)
		assert.False(t, copied.node == document.node, "not the same *html.NewNode")
	})
}

func TestDocument_IsSameNode(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)

		document := &Document{node: parsedDocument}
		require.True(t, document.IsSameNode(document))

		parsedDocument2, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		require.False(t, document.IsSameNode(&Document{node: parsedDocument2}))
	})

	t.Run("same html different address", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
		parsedDocument1, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)
		parsedDocument2, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: parsedDocument1}).IsSameNode(&Document{node: parsedDocument2}))
	})

	t.Run("receiver node is nil", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
		parsedDocument, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: nil}).IsSameNode(&Document{node: parsedDocument}))
	})

	t.Run("param is nil", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
		parsedDocument, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: parsedDocument}).IsSameNode(nil))
	})

	t.Run("param node is nil", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`
		parsedDocument, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: parsedDocument}).IsSameNode(&Document{node: nil}))
	})
}

func TestDocument_GetElementsByTagName(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span id='a'></span><span id='b'></span></body</html>`
	parsedDocument, err := html.Parse(strings.NewReader(textHTML))
	require.NoError(t, err)
	document := &Document{node: parsedDocument}
	elements := document.GetElementsByTagName("span")
	assert.Equal(t, elements.Length(), 2)
}

func TestDocument_GetElementsByClassName(t *testing.T) {
	t.Run("nothing found", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1'></span>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}
		result := document.GetElementsByClassName("nothing-has-this-class")
		require.Nil(t, result)
	})

	t.Run("element found", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1' class='find-me'></span>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}
		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 1, result.Length())
		assert.Equal(t, "x1", result.Item(0).ID())
	})

	t.Run("multiple elements found", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1' class='find-me'></span>
	<span id='x2' class='find-me'></span>
	<div>
		<span id='x3' class='find-me'></span>
	</div>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}

		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 3, result.Length())
		assert.Equal(t, "x1", result.Item(0).ID())
		assert.Equal(t, "x2", result.Item(1).ID())
		assert.Equal(t, "x3", result.Item(2).ID())
	})

	t.Run("embedded element found", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<div id='x1' class='find-me'>
		<div id='x2' class='find-me'>
			<div id='x3' class='find-me'></div>
		</div>
	</div>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}

		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 3, result.Length())
		assert.Equal(t, "x1", result.Item(0).ID())
		assert.Equal(t, "x2", result.Item(1).ID())
		assert.Equal(t, "x3", result.Item(2).ID())
	})

	t.Run("other class in in attribute before", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1' class='other find-me'></span>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}
		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 1, result.Length())
		assert.Equal(t, "x1", result.Item(0).ID())
	})

	t.Run("other class in in attribute after", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1' class='find-me other'></span>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}
		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 1, result.Length())
		assert.Equal(t, "x1", result.Item(0).ID())
	})

	t.Run("extra whitespace", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1' class='    find-me     '></span>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}
		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 1, result.Length())
		assert.Equal(t, "x1", result.Item(0).ID())
	})

	t.Run("multiple classes found", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<span id='x1' class='find-me other'></span>
</body>
</html>`
		parsedDocument, err := html.Parse(strings.NewReader(textHTML))
		require.NoError(t, err)
		document := &Document{node: parsedDocument}
		result := document.GetElementsByClassName("find-me")
		require.Equal(t, 1, result.Length())
	})
}

func TestDocument_QuerySelector(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<div id='1'>
		<span></span>
		<span data-find-me='x' id='x1'></span>
	</div>
</body>
</html>`
	parsedDocument, err := html.Parse(strings.NewReader(textHTML))
	require.NoError(t, err)
	document := &Document{node: parsedDocument}

	result := document.QuerySelector(`#1 [data-find-me="x"]`)
	assert.Equal(t, "x1", result.ID())

	shouldNotBeFound := document.QuerySelector(`#banana-banana-banana`)
	assert.Nil(t, shouldNotBeFound)
}

func TestDocument_QuerySelectorAll(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<div id='1'>
		<span></span>
		<span data-find-me='x' id='x1'></span>
		<span data-find-me='y' id='x2'></span>
	</div>
	<span data-find-me='NOPE' id='x3'></span>
</body>
</html>`
	parsedDocument, err := html.Parse(strings.NewReader(textHTML))
	require.NoError(t, err)
	document := &Document{node: parsedDocument}

	result := document.QuerySelectorAll(`#1 [data-find-me]`)
	require.Equal(t, result.Length(), 2)
	assert.Equal(t, "x1", result.Item(0).ID())
	assert.Equal(t, "x2", result.Item(1).ID())
}

func TestDocument_Contains(t *testing.T) {
	t.Run("a child", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		document := Document{
			node: parsedDocument,
		}
		require.True(t, document.Contains(&Element{node: parsedDocument.LastChild.LastChild}))
	})

	t.Run("some other node", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		document := Document{
			node: parsedDocument,
		}
		require.False(t, document.Contains(&Element{node: &html.Node{}}))
	})

	t.Run("nil", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		document := Document{
			node: parsedDocument,
		}
		require.False(t, document.Contains(&Element{node: nil}))
	})
}

func TestDocument_TextContent(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang='us-en'>
<head><title>Get Elements By Class Name</title></head>
<body>
	<p>Hello, world!</p>
</body>
</html>`
	parsedDocument, err := html.Parse(strings.NewReader(textHTML))
	require.NoError(t, err)
	document := &Document{node: parsedDocument}

	require.Zero(t, document.TextContent())
}

func TestDocument_CreateElement(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.LastChild.FirstChild
	exp.Parent = nil

	document := &Document{node: parsedDocument}
	got := document.CreateElement("SPAN").(*Element).node

	assert.Equal(t, exp, got)
}

func TestDocument_CreateElementIs(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title></title></head><body><div is='fruit'></div></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.LastChild.FirstChild
	exp.Parent = nil

	document := &Document{node: parsedDocument}
	got := document.CreateElementIs("div", "fruit").(*Element).node

	assert.Equal(t, exp, got)
}

func TestDocument_CreateTextNode(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang='us-en'><head><title>peach</title></head><body></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild
	exp.Parent = nil

	document := &Document{node: parsedDocument}
	got := document.CreateTextNode("peach").(*Text).node

	assert.Equal(t, exp, got)
}
