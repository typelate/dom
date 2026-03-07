//go:build js

// Package browser implements the spec interfaces using the browser DOM via
// syscall/js. It is experimental.
package browser

import (
	"iter"
	"syscall/js"

	"github.com/typelate/dom/spec"
)

type Node struct {
	value js.Value
}

func (n *Node) NodeType() spec.NodeType         { return nodeType(n.value) }
func (n *Node) CloneNode(deep bool) spec.Node   { return cloneNode(n.value, deep) }
func (n *Node) IsSameNode(other spec.Node) bool { return isSameNode(n.value, other) }
func (n *Node) TextContent() string             { return textContent(n.value) }

func (n *Node) CompareDocumentPosition(other spec.Node) spec.DocumentPosition {
	return compareDocumentPosition(n.value, other)
}

type Document struct {
	value js.Value
}

func OpenDocument() spec.Document {
	return newDocument(js.Global().Get("document"))
}

func newDocument(value js.Value) spec.Document {
	if value.IsNull() {
		return nil
	}
	return &Document{value: value}
}

func (d *Document) NodeType() spec.NodeType         { return nodeType(d.value) }
func (d *Document) CloneNode(deep bool) spec.Node   { return cloneNode(d.value, deep) }
func (d *Document) IsSameNode(other spec.Node) bool { return isSameNode(d.value, other) }
func (d *Document) TextContent() string             { return textContent(d.value) }

func (d *Document) CompareDocumentPosition(other spec.Node) spec.DocumentPosition {
	return compareDocumentPosition(d.value, other)
}

func (d *Document) Head() spec.Element { return newElement(d.value.Get("head")) }
func (d *Document) Body() spec.Element { return newElement(d.value.Get("body")) }

func (d *Document) Contains(other spec.Node) bool { return contains(d.value, other) }

func (d *Document) GetElementsByTagName(name string) spec.ElementCollection {
	return getElementsByTagName(d.value, name)
}

func (d *Document) GetElementsByClassName(name string) spec.ElementCollection {
	return getElementsByClassName(d.value, name)
}

func (d *Document) QuerySelector(query string) spec.Element {
	return querySelector(d.value, query)
}

func (d *Document) QuerySelectorAll(query string) spec.NodeList[spec.Element] {
	return querySelectorAll(d.value, query)
}

func (d *Document) QuerySelectorSequence(query string) iter.Seq[spec.Element] {
	return querySelectorEach(d.value, query)
}

func (d *Document) CreateElement(localName string) spec.Element {
	return createElement(d.value, localName)
}

func (d *Document) CreateElementIs(localName, is string) spec.Element {
	return createElementIs(d.value, localName, is)
}

func (d *Document) CreateTextNode(text string) spec.Text {
	return createTextNode(d.value, text)
}

type DocumentFragment struct {
	value js.Value
}

func (d *DocumentFragment) NodeType() spec.NodeType         { return nodeType(d.value) }
func (d *DocumentFragment) CloneNode(deep bool) spec.Node   { return cloneNode(d.value, deep) }
func (d *DocumentFragment) IsSameNode(other spec.Node) bool { return isSameNode(d.value, other) }
func (d *DocumentFragment) TextContent() string             { return textContent(d.value) }

func (d *DocumentFragment) CompareDocumentPosition(other spec.Node) spec.DocumentPosition {
	return compareDocumentPosition(d.value, other)
}

func (d *DocumentFragment) Children() spec.ElementCollection { return children(d.value) }
func (d *DocumentFragment) FirstElementChild() spec.Element  { return firstElementChild(d.value) }
func (d *DocumentFragment) LastElementChild() spec.Element   { return lastElementChild(d.value) }
func (d *DocumentFragment) ChildElementCount() int           { return childElementCount(d.value) }

func (d *DocumentFragment) Append(nodes ...spec.Node)          { appendNodes(d.value, nodes) }
func (d *DocumentFragment) Prepend(nodes ...spec.Node)         { prependNodes(d.value, nodes) }
func (d *DocumentFragment) ReplaceChildren(nodes ...spec.Node) { replaceChildrenNodes(d.value, nodes) }

func (d *DocumentFragment) QuerySelector(query string) spec.Element {
	return querySelector(d.value, query)
}

func (d *DocumentFragment) QuerySelectorAll(query string) spec.NodeList[spec.Element] {
	return querySelectorAll(d.value, query)
}

func (d *DocumentFragment) QuerySelectorSequence(query string) iter.Seq[spec.Element] {
	return querySelectorEach(d.value, query)
}

type Element struct {
	value js.Value
}

func newElement(value js.Value) spec.Element {
	if value.IsNull() {
		return nil
	}
	return &Element{value: value}
}

func (e *Element) NodeType() spec.NodeType         { return nodeType(e.value) }
func (e *Element) CloneNode(deep bool) spec.Node   { return cloneNode(e.value, deep) }
func (e *Element) IsSameNode(other spec.Node) bool { return isSameNode(e.value, other) }
func (e *Element) TextContent() string             { return textContent(e.value) }
func (e *Element) CompareDocumentPosition(other spec.Node) spec.DocumentPosition {
	return compareDocumentPosition(e.value, other)
}

func (e *Element) Length() int { return e.Length() }

func (e *Element) IsConnected() bool               { return isConnected(e.value) }
func (e *Element) OwnerDocument() spec.Document    { return ownerDocument(e.value) }
func (e *Element) ParentNode() spec.Node           { return parentNode(e.value) }
func (e *Element) ParentElement() spec.Element     { return parentElement(e.value) }
func (e *Element) PreviousSibling() spec.ChildNode { return previousSibling(e.value) }
func (e *Element) NextSibling() spec.ChildNode     { return nextSibling(e.value) }

func (e *Element) Children() spec.ElementCollection { return children(e.value) }
func (e *Element) FirstElementChild() spec.Element  { return firstElementChild(e.value) }
func (e *Element) LastElementChild() spec.Element   { return lastElementChild(e.value) }
func (e *Element) ChildElementCount() int           { return childElementCount(e.value) }

func (e *Element) Prepend(nodes ...spec.Node) { appendNodes(e.value, nodes) }

func (e *Element) Append(nodes ...spec.Node) { prependNodes(e.value, nodes) }

func (e *Element) ReplaceChildren(nodes ...spec.Node) { replaceChildrenNodes(e.value, nodes) }

func (e *Element) Contains(other spec.Node) bool { return contains(e.value, other) }

func (e *Element) GetElementsByTagName(name string) spec.ElementCollection {
	return getElementsByTagName(e.value, name)
}

func (e *Element) GetElementsByClassName(name string) spec.ElementCollection {
	return getElementsByClassName(e.value, name)
}

func (e *Element) QuerySelector(query string) spec.Element {
	return querySelector(e.value, query)
}

func (e *Element) QuerySelectorAll(query string) spec.NodeList[spec.Element] {
	return querySelectorAll(e.value, query)
}

func (e *Element) QuerySelectorSequence(query string) iter.Seq[spec.Element] {
	return querySelectorEach(e.value, query)
}

func (e *Element) HasChildNodes() bool { return e.value.Bool() }

func (e *Element) ChildNodes() spec.NodeList[spec.Node] {
	nodes := e.value.Get("childNodes")
	if nodes.IsNull() {
		return nil
	}
	return nodeList{value: nodes}
}

func (e *Element) FirstChild() spec.ChildNode { return newChildNode(e.value.Get("firstChild")) }
func (e *Element) LastChild() spec.ChildNode  { return newChildNode(e.value.Get("lastChild")) }

func (e *Element) InsertBefore(node, child spec.ChildNode) spec.ChildNode {
	return newChildNode(e.value.Call("insertBefore", node, child))
}

func (e *Element) AppendChild(node spec.ChildNode) spec.ChildNode {
	return newChildNode(e.value.Call("appendChild", JSValue(node)))
}

func (e *Element) ReplaceChild(node, child spec.ChildNode) spec.ChildNode {
	return newChildNode(e.value.Call("replaceChild", JSValue(node), JSValue(child)))
}

func (e *Element) RemoveChild(node spec.ChildNode) spec.ChildNode {
	return newChildNode(e.value.Call("removeChild", JSValue(node)))
}

func (e *Element) TagName() string   { return e.value.Get("tagName").String() }
func (e *Element) ID() string        { return e.value.Get("id").String() }
func (e *Element) ClassName() string { return e.value.Get("className").String() }

func (e *Element) GetAttribute(name string) string {
	return e.value.Call("getAttribute", name).String()
}
func (e *Element) SetAttribute(name, value string) { e.value.Call("setAttribute", name, value) }

func (e *Element) RemoveAttribute(name string) { e.value.Call("removeAttribute", name) }
func (e *Element) ToggleAttribute(name string) bool {
	return e.value.Call("toggleAttribute", name).Bool()
}
func (e *Element) HasAttribute(name string) bool { return e.value.Call("hasAttribute", name).Bool() }
func (e *Element) Closest(selector string) spec.Element {
	return newElement(e.value.Call("closest", selector))
}
func (e *Element) Matches(selector string) bool { return e.value.Call("matches", selector).Bool() }

func (e *Element) SetInnerHTML(s string) { e.value.Set("innerHTML", s) }
func (e *Element) InnerHTML() string     { return e.value.Get("innerHTML").String() }
func (e *Element) SetOuterHTML(s string) { e.value.Set("innerHTML", s) }
func (e *Element) OuterHTML() string     { return e.value.Get("outerHTML").String() }

type Text struct {
	value js.Value
}

func newTextNode(v js.Value) spec.Text {
	if v.IsNull() {
		return nil
	}
	return &Text{value: v}
}

func (t *Text) NodeType() spec.NodeType         { return nodeType(t.value) }
func (t *Text) CloneNode(deep bool) spec.Node   { return cloneNode(t.value, deep) }
func (t *Text) IsSameNode(other spec.Node) bool { return isSameNode(t.value, other) }
func (t *Text) TextContent() string             { return textContent(t.value) }
func (t *Text) Length() int                     { return t.value.Length() }

func (t *Text) CompareDocumentPosition(other spec.Node) spec.DocumentPosition {
	return compareDocumentPosition(t.value, other)
}

func (t *Text) IsConnected() bool               { return isConnected(t.value) }
func (t *Text) OwnerDocument() spec.Document    { return ownerDocument(t.value) }
func (t *Text) ParentNode() spec.Node           { return parentNode(t.value) }
func (t *Text) ParentElement() spec.Element     { return parentElement(t.value) }
func (t *Text) PreviousSibling() spec.ChildNode { return previousSibling(t.value) }
func (t *Text) NextSibling() spec.ChildNode     { return nextSibling(t.value) }

func (t *Text) Data() string     { return t.value.Get("data").String() }
func (t *Text) SetData(s string) { t.value.Set("data", s) }

var (
	nodeClass             = js.Global().Get("Node")
	textClass             = js.Global().Get("Text")
	documentClass         = js.Global().Get("Document")
	documentFragmentClass = js.Global().Get("DocumentFragment")
	elementClass          = js.Global().Get("Element")
)

func NewNode(value js.Value) spec.Node {
	if value.InstanceOf(elementClass) {
		return newElement(value)
	}
	if value.InstanceOf(textClass) {
		return newTextNode(value)
	}
	if value.InstanceOf(documentClass) {
		return newDocument(value)
	}
	if value.InstanceOf(documentFragmentClass) {
		return &DocumentFragment{value: value}
	}
	if value.InstanceOf(nodeClass) {
		return &Node{value: value}
	}
	return nil
}

func newChildNode(value js.Value) spec.ChildNode {
	cn, ok := NewNode(value).(spec.ChildNode)
	if !ok {
		return nil
	}
	return cn
}

func JSValue(node any) js.Value {
	switch n := node.(type) {
	case *Element:
		return n.value
	case *Node:
		return n.value
	case *Document:
		return n.value
	case *DocumentFragment:
		return n.value
	case *Text:
		return n.value
	case js.Value:
		return n
	default:
		return js.Null()
	}
}

func valueArray(in []spec.Node) []any {
	out := make([]any, 0, len(in))
	for _, n := range in {
		out = append(out, JSValue(n))
	}
	return out
}

func nodeType(receiver js.Value) spec.NodeType {
	return spec.NodeType(receiver.Get("nodeType").Int())
}

func cloneNode(receiver js.Value, deep bool) spec.Node {
	return NewNode(receiver.Call("cloneNode", deep))
}

func isSameNode(receiver js.Value, other spec.Node) bool {
	return receiver.Call("isSameNode", JSValue(other)).Bool()
}

func textContent(receiver js.Value) string {
	return receiver.Get("textContent").String()
}

func contains(receiver js.Value, other spec.Node) bool {
	return receiver.Call("contains", JSValue(other)).Bool()
}

func item(receiver js.Value, index int) spec.Element {
	return newElement(receiver.Call("item", index))
}

type htmlCollection struct {
	value js.Value
}

func (e htmlCollection) Length() int { return e.value.Length() }

func (e htmlCollection) Item(index int) spec.Element {
	return item(e.value, index)
}

func (e htmlCollection) NamedItem(name string) spec.Element {
	return newElement(e.value.Call("namedItem", name))
}

func getElementsByTagName(receiver js.Value, name string) spec.ElementCollection {
	return htmlCollection{value: receiver.Call("getElementsByTagName", name)}
}

func getElementsByClassName(receiver js.Value, name string) spec.ElementCollection {
	return htmlCollection{value: receiver.Call("getElementsByClassName", name)}
}

func querySelector(receiver js.Value, query string) spec.Element {
	return newElement(receiver.Call("querySelector", query))
}

type elementList struct {
	value js.Value
}

func (e elementList) Length() int { return e.value.Length() }

func (e elementList) Item(index int) spec.Element {
	return item(e.value, index)
}

func querySelectorAll(receiver js.Value, query string) elementList {
	return elementList{value: receiver.Call("querySelectorAll", query)}
}

func querySelectorEach(receiver js.Value, query string) iter.Seq[spec.Element] {
	return func(f func(spec.Element) bool) {
		list := querySelectorAll(receiver, query)
		for i := 0; i < list.Length(); i++ {
			if !f(list.Item(i)) {
				return
			}
		}
	}
}

func createElement(receiver js.Value, tagName string) spec.Element {
	return newElement(receiver.Call("createElement", tagName))
}

func createElementIs(receiver js.Value, tagName, is string) spec.Element {
	return newElement(receiver.Call("createElement", tagName, struct {
		Is string `json:"is"`
	}{Is: is}))
}

func createTextNode(receiver js.Value, text string) spec.Text {
	n := receiver.Call("createTextNode", text)
	if n.IsNull() {
		return nil
	}
	return &Text{value: n}
}

func isConnected(receiver js.Value) bool { return receiver.Get("isConnected").Bool() }
func ownerDocument(receiver js.Value) spec.Document {
	return newDocument(receiver.Get("ownerDocument"))
}
func parentNode(receiver js.Value) spec.Node       { return NewNode(receiver.Get("parentNode")) }
func parentElement(receiver js.Value) spec.Element { return newElement(receiver.Get("parentElement")) }
func children(receiver js.Value) htmlCollection {
	return htmlCollection{value: receiver.Call("children")}
}

func compareDocumentPosition(receiver js.Value, other spec.Node) spec.DocumentPosition {
	return spec.DocumentPosition(receiver.Call("compareDocumentPosition", JSValue(other)).Int())
}

func previousSibling(receiver js.Value) spec.ChildNode {
	return newChildNode(receiver.Get("previousSibling"))
}

func nextSibling(receiver js.Value) spec.ChildNode {
	return newChildNode(receiver.Get("nextSibling"))
}

func firstElementChild(receiver js.Value) spec.Element {
	return newElement(receiver.Get("firstElementChild"))
}

func lastElementChild(receiver js.Value) spec.Element {
	return newElement(receiver.Get("lastElementChild"))
}
func childElementCount(receiver js.Value) int { return receiver.Get("childElementCount").Int() }

func appendNodes(receiver js.Value, in []spec.Node) {
	receiver.Call("append", valueArray(in)...)
}

func prependNodes(receiver js.Value, in []spec.Node) {
	receiver.Call("prepend", valueArray(in)...)
}

func replaceChildrenNodes(receiver js.Value, in []spec.Node) {
	receiver.Call("replaceChildren", valueArray(in)...)
}

type nodeList struct {
	value js.Value
}

func (n nodeList) Length() int          { return n.value.Length() }
func (n nodeList) Item(i int) spec.Node { return NewNode(n.value.Call("item", i)) }
