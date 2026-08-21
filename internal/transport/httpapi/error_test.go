package httpapi

import (
	"bytes"
	"errors"
	"example.com/api-traffic-governance/internal/domain"
	"fmt"
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

func TestWriteErrRecognizesWrappedInvalidInput(t *testing.T) {
	rr := httptest.NewRecorder()
	writeErr(rr, fmt.Errorf("decode: %w", domain.ErrInvalidInput))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("wrapped input code = %d", rr.Code)
	}
}

func TestWriteErrRecognizesWrappedConflict(t *testing.T) {
	rr := httptest.NewRecorder()
	writeErr(rr, fmt.Errorf("publish: %w", domain.ErrConflict))
	if rr.Code != http.StatusConflict {
		t.Fatalf("wrapped conflict code = %d", rr.Code)
	}
}

type wrappedNotFound struct{}

func (wrappedNotFound) Error() string { return "lookup failed" }
func (wrappedNotFound) Unwrap() error { return domain.ErrNotFound }
