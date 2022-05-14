// (c) Gareth Watts 2017
// Licensed under an MIT license
// See the LICENSE file for details

package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testMiddleware struct {
	h http.Handler
}

func newTestMiddleware(h http.Handler) http.Handler {
	return &testMiddleware{h}
}

type newTestMiddlewareWriter struct {
	http.ResponseWriter
}

func (ntmw newTestMiddlewareWriter) Write(b []byte) (int, error) {
	ntmw.Header().Set("middle-write", "true")
	return ntmw.ResponseWriter.Write(b)
}

func (h *testMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "test", "value")
	w.Header().Set("middle-header", "true")
	h.h.ServeHTTP(newTestMiddlewareWriter{w}, r.WithContext(ctx))
}

func TestContextPassthrough(t *testing.T) {
	assert := assert.New(t)
	engine := gin.New()

	// create some outer middleware
	engine.Use(func(c *gin.Context) {
		assert.Nil(c.Request.Context().Value("test"), "context should not propagate upstream")
		c.Header("outer-pre", "true")
		c.Next()
		c.Header("outer-post", "true")
		assert.Nil(c.Request.Context().Value("test"), "context should not propagate upstream")
	})

	// hook up our standard middleware
	h, adapter := New()
	middle := &testMiddleware{h}
	engine.Use(adapter(middle))

	engine.GET("/", func(c *gin.Context) {
		assert.Equal("value", c.Request.Context().Value("test"), "context should be passsed through from middleware")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal("true", w.Header().Get("outer-pre"))
	assert.Equal("true", w.Header().Get("outer-post"))
	assert.Equal("true", w.Header().Get("middle-header"))
}

func TestAbort(t *testing.T) {
	assert := assert.New(t)
	engine := gin.New()

	// Create middleware that does not call its wrapped handler
	middle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("middle-header", "true")
	})

	_, adapter := New()
	engine.Use(adapter(middle))

	engine.GET("/", func(c *gin.Context) {
		assert.Fail("Should not be reached")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal("true", w.Header().Get("middle-header"))

}

func TestWrap(t *testing.T) {
	assert := assert.New(t)
	engine := gin.New()

	engine.Use(Wrap(newTestMiddleware))
	engine.GET("/", func(c *gin.Context) {
		assert.Equal("value", c.Request.Context().Value("test"), "context should be passsed through from middleware")

		// This Write validates that the writer, modified by the wrapped middleware, is propagated.
		c.Writer.Write([]byte("test"))
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal("true", w.Header().Get("middle-header"))
	assert.Equal("true", w.Header().Get("middle-write"))
}
