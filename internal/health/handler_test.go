package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeDependencyChecker struct {
	err error
}

func (f *fakeDependencyChecker) Check() error {
	return f.err
}

func TestLiveHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	recorder := httptest.NewRecorder()

	LiveHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestReadyHandler_WhenDependencyIsHealthy(t *testing.T) {
	checker := &fakeDependencyChecker{}
	handler := NewHandler(checker)

	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	recorder := httptest.NewRecorder()

	handler.ReadyHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestReadyHandler_WhenDependencyIsUnavailable(t *testing.T) {
	checker := &fakeDependencyChecker{
		err: assertError{},
	}
	handler := NewHandler(checker)

	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	recorder := httptest.NewRecorder()

	handler.ReadyHandler(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			recorder.Code,
		)
	}
}

type assertError struct{}

func (assertError) Error() string {
	return "dependency unavailable"
}
