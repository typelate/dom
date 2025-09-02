//go:build js

package browser_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/typelate/dom/browser"
	"github.com/typelate/dom/spec"
)

func TestDocument(t *testing.T) {
	const childClass = "child"

	type scope struct {
		document spec.Document
		head,
		body,
		div,
		childOne,
		childTwo spec.Element
	}

	setup := sync.OnceValue(func() scope {
		document := browser.OpenDocument()
		div := document.CreateElement("div")

		head := document.Head()
		require.NotNil(t, head)

		body := document.Body()
		require.NotNil(t, body)

		childOne := document.CreateElement("div")
		childOne.SetAttribute("id", "one")
		childOne.SetAttribute("class", childClass)
		assert.NotNil(t, body.AppendChild(childOne))

		childTwo := document.CreateElement("div")
		childTwo.SetAttribute("id", "two")
		childTwo.SetAttribute("class", childClass)
		assert.NotNil(t, body.AppendChild(childTwo))

		return scope{
			document: document,
			head:     head,
			body:     body,
			div:      div,
			childOne: childOne,
			childTwo: childTwo,
		}
	})

	t.Run("NodeType", func(t *testing.T) {
		ts := setup()
		assert.Equal(t, spec.NodeTypeDocument, ts.document.NodeType())
	})
	t.Run("CloneNode", func(t *testing.T) {
		ts := setup()
		assert.NotPanics(t, func() { ts.document.CloneNode(false) })
		assert.NotPanics(t, func() { ts.document.CloneNode(true) })
	})
	t.Run("IsSameNode", func(t *testing.T) {
		ts := setup()
		assert.True(t, ts.document.IsSameNode(ts.document))
		assert.False(t, ts.document.IsSameNode(ts.div))
	})
	t.Run("TextContent", func(t *testing.T) {
		ts := setup()
		assert.NotPanics(t, func() { ts.document.TextContent() })
	})
	t.Run("Contains", func(t *testing.T) {
		ts := setup()
		assert.True(t, ts.document.Contains(ts.head))
		assert.False(t, ts.document.Contains(ts.div))
	})
	t.Run("GetElementsByTagName", func(t *testing.T) {
		ts := setup()
		assert.NotZero(t, ts.document.GetElementsByTagName("div").Length())
	})
	t.Run("GetElementsByClassName", func(t *testing.T) {
		ts := setup()
		assert.Equal(t, 2, ts.document.GetElementsByClassName(childClass).Length())
	})
}

func TestElement_CompareDocumentPosition(t *testing.T) {
	type scope struct {
		document spec.Document
		a, b, c  spec.Element
	}

	setup := sync.OnceValue(func() scope {
		document := browser.OpenDocument()

		a := document.CreateElement("div")
		a.SetAttribute("id", "a")

		b := document.CreateElement("span")
		b.SetAttribute("id", "b")
		a.AppendChild(b)

		c := document.CreateElement("div")
		c.SetAttribute("id", "c")

		document.Body().AppendChild(a)
		document.Body().AppendChild(c)

		return scope{
			a:        a,
			b:        b,
			c:        c,
			document: document,
		}
	})

	t.Run("same node", func(t *testing.T) {
		ts := setup()
		assert.Equal(t, spec.DocumentPosition(0), ts.a.CompareDocumentPosition(ts.a))
	})

	t.Run("contains", func(t *testing.T) {
		ts := setup()
		pos := ts.a.CompareDocumentPosition(ts.b)
		assert.Equal(t, spec.DocumentPositionContainedBy|spec.DocumentPositionFollowing, pos)
	})

	t.Run("contained by", func(t *testing.T) {
		ts := setup()
		pos := ts.b.CompareDocumentPosition(ts.a)
		assert.Equal(t, spec.DocumentPositionContains|spec.DocumentPositionPreceding, pos)
	})

	t.Run("preceding", func(t *testing.T) {
		ts := setup()
		pos := ts.a.CompareDocumentPosition(ts.c)
		assert.Equal(t, spec.DocumentPositionFollowing, pos)
	})

	t.Run("following", func(t *testing.T) {
		ts := setup()
		pos := ts.c.CompareDocumentPosition(ts.a)
		assert.Equal(t, spec.DocumentPositionPreceding, pos)
	})

	t.Run("disconnected", func(t *testing.T) {
		ts := setup()
		d := ts.document.CreateElement("div")
		// not appended to the DOM
		pos := ts.a.CompareDocumentPosition(d)
		assert.True(t, spec.DocumentPositionDisconnected&pos != 0)
		assert.True(t, spec.DocumentPositionImplementationSpecific&pos != 0)
	})
}
