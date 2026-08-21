package httpapi

import (
	"bytes"
	"errors"
	"example.com/api-traffic-governance/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteErrRecognizesWrappedSentinel(t *testing.T) {
	rr := httptest.NewRecorder()
	writeErr(rr, errors.New("outer: "+domain.ErrNotFound.Error()))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("plain text control code = %d", rr.Code)
	}
	rr = httptest.NewRecorder()
	writeErr(rr, wrappedNotFound{})
	if rr.Code != http.StatusNotFound {
		t.Fatalf("wrapped not found code = %d", rr.Code)
	}
	_ = bytes.NewBuffer(nil)
}

type wrappedNotFound struct{}

func (wrappedNotFound) Error() string { return "lookup failed" }
func (wrappedNotFound) Unwrap() error { return domain.ErrNotFound }
