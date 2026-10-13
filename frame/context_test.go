package frame

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Fail 应把错误写成一个 JSON 响应
func TestContextFail_WritesJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	c := newContext(w, httptest.NewRequest("GET", "/", nil))

	c.Fail(http.StatusInternalServerError, "boom")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not valid JSON: %v (body=%q)", err, w.Body.String())
	}
	if body["message"] != "boom" {
		t.Errorf("message = %q, want %q", body["message"], "boom")
	}
}

// Fail 之后，链上剩余的 handler 不应再执行
func TestContextFail_AbortsRemainingChain(t *testing.T) {
	w := httptest.NewRecorder()
	c := newContext(w, httptest.NewRequest("GET", "/", nil))

	reached := false
	c.handlers = []HandlerFunc{
		func(c *Context) { c.Fail(http.StatusInternalServerError, "boom") },
		func(c *Context) { reached = true },
	}

	c.Next()

	if reached {
		t.Error("Fail 之后的 handler 仍然执行了，链没有被中断")
	}
}
