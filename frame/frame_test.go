package frame

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

// 静态文件命中时，应把请求交给 fileServer 处理
func TestCreateStaticHandler_ServesExistingFile(t *testing.T) {
	fsys := fstest.MapFS{
		"css/style.css": &fstest.MapFile{Data: []byte("body{color:red}")},
	}

	e := New()
	handler := e.createStaticHandler("/assets", http.FS(fsys))

	w := httptest.NewRecorder()
	c := newContext(w, httptest.NewRequest("GET", "/assets/css/style.css", nil))
	c.Params = map[string]string{"filepath": "css/style.css"}

	handler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if body := w.Body.String(); !strings.Contains(body, "color:red") {
		t.Errorf("body = %q, want it to contain %q", body, "color:red")
	}
}

// 文件不存在时应返回 404
func TestCreateStaticHandler_MissingFileIsNotFound(t *testing.T) {
	e := New()
	handler := e.createStaticHandler("/assets", http.FS(fstest.MapFS{}))

	w := httptest.NewRecorder()
	c := newContext(w, httptest.NewRequest("GET", "/assets/nope.txt", nil))
	c.Params = map[string]string{"filepath": "nope.txt"}

	handler(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
