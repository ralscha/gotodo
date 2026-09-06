package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	if ok := JSON(recorder, http.StatusCreated, map[string]string{"status": "ok"}); !ok {
		t.Fatal("JSON() = false, want true")
	}
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := recorder.Body.String(); got != `{"status":"ok"}` {
		t.Fatalf("body = %q", got)
	}
}

func TestJSONMarshalFailure(t *testing.T) {
	recorder := httptest.NewRecorder()

	if ok := JSON(recorder, http.StatusOK, make(chan int)); ok {
		t.Fatal("JSON() = true, want false")
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}

func TestJSONWriteFailureDoesNotWriteSecondResponse(t *testing.T) {
	w := &failingResponseWriter{header: make(http.Header)}

	if ok := JSON(w, http.StatusOK, map[string]bool{"ok": true}); ok {
		t.Fatal("JSON() = true, want false")
	}
	if w.writeHeaderCalls != 1 {
		t.Fatalf("WriteHeader called %d times, want 1", w.writeHeaderCalls)
	}
}

type failingResponseWriter struct {
	header           http.Header
	writeHeaderCalls int
}

func (w *failingResponseWriter) Header() http.Header { return w.header }

func (w *failingResponseWriter) WriteHeader(int) { w.writeHeaderCalls++ }

func (w *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}
