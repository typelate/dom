package dom_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom"
	"github.com/typelate/dom/spec"
)

func parseDocumentFragment(t *testing.T, input string) *dom.DocumentFragment {
	t.Helper()
	children, err := html.ParseFragment(strings.NewReader(input), &html.Node{
		Type:     html.ElementNode,
		Data:     atom.Div.String(),
		DataAtom: atom.Div,
	})
	require.NoError(t, err)
	return dom.NewDocumentFragment(children)
}

func TestDocumentFragment_NodeType(t *testing.T) {
	fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
	require.Equal(t, spec.NodeTypeDocumentFragment, fragment.NodeType())
}

func TestDocumentFragment_CloneNode(t *testing.T) {
	t.Run("deep", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)

		clone := fragment.CloneNode(true)
		clonedFragment := clone.(*dom.DocumentFragment)

		require.False(t, fragment.IsSameNode(clone))
		children := fragment.Children()
		clonedChildren := clonedFragment.Children()
		for i := 0; i < children.Length(); i++ {
			require.False(t, clonedChildren.Item(i).IsSameNode(children.Item(i)))
		}
	})
	t.Run("shallow", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)

		clone := fragment.CloneNode(false)
		clonedFragment := clone.(*dom.DocumentFragment)

		require.False(t, fragment.IsSameNode(clone))
		children := fragment.Children()
		clonedChildren := clonedFragment.Children()
		for i := 0; i < children.Length(); i++ {
			require.True(t, clonedChildren.Item(i).IsSameNode(children.Item(i)))
		}
	})
}

func TestDocumentFragment_IsSameNode(t *testing.T) {
	t.Run("another fragment", func(t *testing.T) {
		a := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
		require.True(t, a.IsSameNode(a))

		b := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)

		require.False(t, a.IsSameNode(b))
		require.False(t, b.IsSameNode(a))

		c := b
		require.True(t, b.IsSameNode(c))
		require.True(t, c.IsSameNode(b))
		require.True(t, c.IsSameNode(c))
	})
	t.Run("another node type", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
		require.False(t, fragment.IsSameNode(nil))
	})
}

func TestDocumentFragment_TextContent(t *testing.T) {
	fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
	require.Equal(t, "Hello, world!", fragment.TextContent())
}

func TestDocumentFragment_Children(t *testing.T) {
	fragment := parseDocumentFragment(t, `<div>1</div><div>2</div><div>3</div>`)
	children := fragment.Children()
	require.Equal(t, 3, children.Length())
	require.Equal(t, "1", children.Item(0).TextContent())
	require.Equal(t, "2", children.Item(1).TextContent())
	require.Equal(t, "3", children.Item(2).TextContent())
}

func TestDocumentFragment_FirstElementChild(t *testing.T) {
	t.Run("has elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `a<span>x</span>b<span>y</span>c`)
		require.Equal(t, "x", fragment.FirstElementChild().TextContent())
	})
	t.Run("no elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello!`)
		require.Nil(t, fragment.FirstElementChild())
	})
	t.Run("one element", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `hello<div>x</div>`)
		require.Equal(t, "x", fragment.FirstElementChild().TextContent())
	})
}

func TestDocumentFragment_LastElementChild(t *testing.T) {
	t.Run("has elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `a<span>x</span>b<span>y</span>c`)
		require.Equal(t, "y", fragment.LastElementChild().TextContent())
	})
	t.Run("no elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello!`)
		require.Nil(t, fragment.LastElementChild())
	})
	t.Run("one element", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div>x</div>world`)
		require.Equal(t, "x", fragment.LastElementChild().TextContent())
	})
}

func TestDocumentFragment_ChildElementCount(t *testing.T) {
	fragment := parseDocumentFragment(t, `a<div>x</div>b<div>y</div>`)
	require.Equal(t, 2, fragment.ChildElementCount())
}

func TestDocumentFragment_Append(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fragment := parseDocumentFragment(t, ``)
		fragment.Append()
		require.Equal(t, 0, fragment.ChildElementCount())
	})
	t.Run("add one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 1)
		node := dom.NewNode(other[0]).(spec.Element)

		fragment.Append(node)

		require.Equal(t, 2, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "a", children.Item(0).ID())
		assert.Equal(t, "b", children.Item(1).ID())
	})
	t.Run("add two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div><div id="c"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 2)

		a, ok := dom.NewNode(other[0]).(spec.Element)
		require.True(t, ok)

		b, ok := dom.NewNode(other[1]).(spec.Element)
		require.True(t, ok)

		fragment.Append(a, b)

		require.Equal(t, 3, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "a", children.Item(0).ID())
		assert.Equal(t, "b", children.Item(1).ID())
		assert.Equal(t, "c", children.Item(2).ID())
	})
}

func TestDocumentFragment_Prepend(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fragment := parseDocumentFragment(t, ``)
		fragment.Prepend()
		require.Equal(t, 0, fragment.ChildElementCount())
	})
	t.Run("add one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 1)
		node := dom.NewNode(other[0]).(spec.Element)

		fragment.Prepend(node)

		require.Equal(t, 2, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
		assert.Equal(t, "a", children.Item(1).ID())
	})
	t.Run("add two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div><div id="c"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 2)

		a, ok := dom.NewNode(other[0]).(spec.Element)
		require.True(t, ok)

		b, ok := dom.NewNode(other[1]).(spec.Element)
		require.True(t, ok)

		fragment.Prepend(a, b)

		require.Equal(t, 3, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
		assert.Equal(t, "c", children.Item(1).ID())
		assert.Equal(t, "a", children.Item(2).ID())
	})
}

func TestDocumentFragment_ReplaceChildren(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fragment := parseDocumentFragment(t, ``)
		fragment.ReplaceChildren()
		require.Equal(t, 0, fragment.ChildElementCount())
	})
	t.Run("replace one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 1)
		node := dom.NewNode(other[0]).(spec.Element)

		fragment.ReplaceChildren(node)

		require.Equal(t, 1, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
	})
	t.Run("replace two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div><div id="c"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 2)

		a, ok := dom.NewNode(other[0]).(spec.Element)
		require.True(t, ok)

		b, ok := dom.NewNode(other[1]).(spec.Element)
		require.True(t, ok)

		fragment.ReplaceChildren(a, b)

		require.Equal(t, 2, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
		assert.Equal(t, "c", children.Item(1).ID())
	})
}

func TestDocumentFragment_QuerySelector(t *testing.T) {
	t.Run("found one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.NotNil(t, fragment.QuerySelector("#a"))
	})
	t.Run("found two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `
<section><div class="a" id="i1"></div></section>
<section></section>
<section><div class="a" id="i2"></div><div class="a" id="i3"></div></section>`)
		result := fragment.QuerySelector(".a")
		require.NotNil(t, result)
		assert.Equal(t, "i1", result.ID())
	})
	t.Run("not found", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.Nil(t, fragment.QuerySelector("#not-found"))
	})
	t.Run("direct descendant", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<p>Peach</p>`)
		result := fragment.QuerySelector("p")
		require.NotNil(t, result)
		assert.Equal(t, "Peach", result.TextContent())
	})
	t.Run("tree of nested results", func(t *testing.T) {
		fragment := parseDocumentFragment(t,
			/* language=html */ `<div id='n0'>
	<div id='n1'>
		<div id='n2'></div>
	</div>
	<div id='n3'></div>
</div>
<div id='n4'>
	<div id='n5'>
		<div id='n6'></div>
	</div>
	<div id='n7'></div>
</div>`)
		results := fragment.QuerySelectorAll("div")
		require.NotNil(t, results)
		assert.Equal(t, results.Length(), 8)

		for i := 0; i < results.Length(); i++ {
			result := results.Item(i)
			require.NotNil(t, result)
			assert.Equal(t, "DIV", result.TagName())
			assert.Equal(t, "n"+strconv.Itoa(i), result.ID())
		}
	})
}

func TestDocumentFragment_QuerySelectorAll(t *testing.T) {
	t.Run("found one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.NotNil(t, fragment.QuerySelectorAll("#a"))
	})
	t.Run("found two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `
<section><div class="a" id="i1"></div></section>
<section></section>
<section><div class="a" id="i2"></div><div class="a" id="i3"></div></section>`)
		results := fragment.QuerySelectorAll(".a")
		assert.Equal(t, "i1", results.Item(0).ID())
		assert.Equal(t, "i2", results.Item(1).ID())
		assert.Equal(t, "i3", results.Item(2).ID())
	})
	t.Run("not found", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.Nil(t, fragment.QuerySelectorAll("#not-found"))
	})
}

func TestDocumentFragment_String(t *testing.T) {
	t.Run("found two", func(t *testing.T) {
		rawHTML := `
<section><div class="a" id="i1"></div></section>
<section></section>
<section><div class="a" id="i2"></div><div class="a" id="i3"></div></section>`

		fragment := parseDocumentFragment(t, rawHTML)

		assert.Equal(t, rawHTML, fragment.String())
	})
}

func TestDocumentFragment_QuerySelectorEach(t *testing.T) {
	fragment, _ := html.ParseFragment(strings.NewReader( /* language=html */ `<div id="n1">
	<div id="n2">
		<div id="n3"></div>
	</div>
	<div id="n4"></div>
</div>
<div id="n5"></div>
<div id="n6">
	<div id="n7"></div>
</div>`), &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body})

	for stopOnCall := 1; stopOnCall <= 7; stopOnCall++ {
		callCount := 0
		dom.NewDocumentFragment(fragment).QuerySelectorSequence("[id]")(func(el spec.Element) bool {
			callCount++
			assert.NotNil(t, el)
			if callCount == stopOnCall {
				return false
			} else {
				return true
			}
		})
		assert.Equal(t, stopOnCall, callCount)
	}
}

// DocumentFragment must satisfy the same query interface as Document and
// Element, so one test helper can accept any of them.
var _ spec.ElementQueries = (spec.DocumentFragment)(nil)

func TestDocumentFragment_Contains(t *testing.T) {
	fragment := parseDocumentFragment(t, `<p class="a">one</p><div><em>two</em></div>`)

	t.Run("a top level element", func(t *testing.T) {
		p := fragment.QuerySelector("p")
		require.NotNil(t, p)
		assert.True(t, fragment.Contains(p))
	})

	t.Run("a nested element", func(t *testing.T) {
		em := fragment.QuerySelector("em")
		require.NotNil(t, em)
		assert.True(t, fragment.Contains(em))
	})

	t.Run("an element from somewhere else", func(t *testing.T) {
		other := parseDocumentFragment(t, `<span>elsewhere</span>`).QuerySelector("span")
		require.NotNil(t, other)
		assert.False(t, fragment.Contains(other))
	})

	t.Run("nil", func(t *testing.T) {
		assert.False(t, fragment.Contains(nil))
	})
}

func TestDocumentFragment_GetElementsByTagName(t *testing.T) {
	fragment := parseDocumentFragment(t, `<p>one</p><div><p>two</p></div>`)

	// The top level nodes are part of the fragment, so a match on one of them
	// counts just as a nested match does.
	list := fragment.GetElementsByTagName("p")
	require.Equal(t, 2, list.Length())
	assert.Equal(t, "one", list.Item(0).TextContent())
	assert.Equal(t, "two", list.Item(1).TextContent())

	assert.Equal(t, 0, fragment.GetElementsByTagName("table").Length())
}

func TestDocumentFragment_GetElementsByClassName(t *testing.T) {
	fragment := parseDocumentFragment(t, `<p class="a b">one</p><div><p class="a">two</p></div>`)

	both := fragment.GetElementsByClassName("a")
	require.Equal(t, 2, both.Length())
	assert.Equal(t, "one", both.Item(0).TextContent())

	only := fragment.GetElementsByClassName("b")
	require.Equal(t, 1, only.Length())
	assert.Equal(t, "one", only.Item(0).TextContent())

	assert.Equal(t, 0, fragment.GetElementsByClassName("missing").Length())
}
