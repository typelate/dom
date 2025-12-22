package domtest_test

import (
	_ "embed"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"

	"github.com/typelate/dom/domtest"
	"github.com/typelate/dom/spec"
)

var (
	//go:embed testdata/index.html
	indexHTML string

	//go:embed testdata/fragment.html
	fragmentHTML string
)

func TestParseResponseDocument(t *testing.T) {
	t.Run("when a valid html document is passed", func(t *testing.T) {
		testingT := newTestingT()
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(indexHTML)),
		}
		document := domtest.ParseResponseDocument(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())

		require.NotNil(t, document)
		p := document.QuerySelector(`p`)
		assert.Equal(t, p.TextContent(), "Hello, world!")
	})

	t.Run("when a fragment is returned", func(t *testing.T) {
		testingT := newTestingT()
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(fragmentHTML)),
		}
		document := domtest.ParseResponseDocument(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())

		require.NotNil(t, document)
		list := document.QuerySelectorAll(`p`)
		require.Equal(t, 2, list.Length())
		require.Equal(t, list.Item(0).TextContent(), "Hello, world!")
		require.Equal(t, list.Item(1).TextContent(), "Greeting!")
	})

	t.Run("when read fails and close is ok", func(t *testing.T) {
		testingT := newTestingT()
		fakeBody := &errClose{
			Reader:   iotest.ErrReader(errors.New("banana")),
			closeErr: nil,
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.ParseResponseDocument(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 1, "it should report an error")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})

	t.Run("when read is ok but close fails", func(t *testing.T) {
		testingT := newTestingT()
		fakeBody := &errClose{
			Reader:   strings.NewReader(indexHTML),
			closeErr: errors.New("banana"),
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.ParseResponseDocument(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 1, "it should report two errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})

	t.Run("when both read and close fail", func(t *testing.T) {
		testingT := newTestingT()
		fakeBody := &errClose{
			Reader:   iotest.ErrReader(errors.New("banana")),
			closeErr: errors.New("lemon"),
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.ParseResponseDocument(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 2, "it should report two errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.Equal(t, testingT.HelperCallCount(), 1)
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})

	t.Run("when both read and close fail", func(t *testing.T) {
		testingT := newTestingT()
		fakeBody := &errClose{
			Reader:   iotest.ErrReader(errors.New("banana")),
			closeErr: errors.New("lemon"),
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.ParseResponseDocument(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 2, "it should report two errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})
}

func TestParseResponseDocumentFragment(t *testing.T) {
	t.Run("when a valid html document is passed", func(t *testing.T) {
		testingT := newTestingT()
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(fragmentHTML)),
		}
		fragment := domtest.ParseResponseDocumentFragment(testingT, res, atom.Body)

		assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())

		require.NotNil(t, fragment)
		require.Equal(t, spec.NodeTypeDocumentFragment, fragment.NodeType())
	})
}

func TestParseReaderDocument(t *testing.T) {
	testingT := newTestingT()
	r := iotest.ErrReader(errors.New("banana"))

	document := domtest.ParseReaderDocument(testingT, r)

	assert.Equal(t, testingT.ErrorCallCount(), 1, "it should report two errors")
	assert.NotZero(t, testingT.HelperCallCount())
	assert.Nil(t, document)
}

func TestParseStringDocument(t *testing.T) {
	testingT := newTestingT()

	document := domtest.ParseStringDocument(testingT, "<p>Hello, world!</p>")

	assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
	assert.Equal(t, testingT.LogCallCount(), 0)
	assert.NotZero(t, testingT.HelperCallCount())

	assert.NotNil(t, document)
	p := document.QuerySelector(`p`)
	assert.Equal(t, p.TextContent(), "Hello, world!")
}

type errClose struct {
	io.Reader
	closeCallCount int
	closeErr       error
}

func (e *errClose) Close() error {
	e.closeCallCount++
	return e.closeErr
}

type TestingT struct {
	mock.Mock
}

func newTestingT() *TestingT {
	t := &TestingT{}
	t.On("Helper").Maybe()
	t.On("Error", mock.Anything).Maybe()
	t.On("Log", mock.Anything).Maybe()
	t.On("Errorf", mock.Anything, mock.Anything).Maybe()
	t.On("FailNow").Maybe()
	t.On("Failed").Return(false).Maybe()
	t.On("SkipNow").Maybe()
	return t
}

func (m *TestingT) Helper() {
	m.Called()
}

func (m *TestingT) Error(args ...any) {
	m.Called(args)
}

func (m *TestingT) Log(args ...any) {
	m.Called(args)
}

func (m *TestingT) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *TestingT) FailNow() {
	m.Called()
}

func (m *TestingT) Failed() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *TestingT) SkipNow() {
	m.Called()
}

// Helper methods for counting calls (compatible with existing tests)

func (m *TestingT) ErrorCallCount() int {
	return m.countCalls("Error")
}

func (m *TestingT) LogCallCount() int {
	return m.countCalls("Log")
}

func (m *TestingT) HelperCallCount() int {
	return m.countCalls("Helper")
}

func (m *TestingT) ErrorfCallCount() int {
	return m.countCalls("Errorf")
}

func (m *TestingT) FailNowCallCount() int {
	return m.countCalls("FailNow")
}

func (m *TestingT) FailedCallCount() int {
	return m.countCalls("Failed")
}

func (m *TestingT) SkipNowCallCount() int {
	return m.countCalls("SkipNow")
}

func (m *TestingT) countCalls(methodName string) int {
	count := 0
	for _, call := range m.Calls {
		if call.Method == methodName {
			count++
		}
	}
	return count
}

var _ domtest.TestingT = new(TestingT)
