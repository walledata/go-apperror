package apperror

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func newRemoteErrorFixture() *RemoteError {
	return &RemoteError{
		Canonical:   NewUnavailable("user-service.GetUser", WithMessage("user-service call failed")),
		Service:     "user-service",
		Operation:   "GetUser",
		Request:     &Request{Method: "GET", URL: "/users/42"},
		Response:    &Response{StatusCode: 503, Body: []byte(`{"code":"DEGRADED"}`)},
		BodyCode:    "DEGRADED",
		BodyMessage: "service degraded",
		RetryAfter:  30 * time.Second,
	}
}

func TestRemoteErrorCanonicalAccessors(t *testing.T) {
	r := newRemoteErrorFixture()
	if r.Canonical.Code() != CodeUnavailable {
		t.Errorf("Canonical.Code() = %s, want %s", r.Canonical.Code(), CodeUnavailable)
	}
	if r.Canonical.Message() != "user-service call failed" {
		t.Errorf("Canonical.Message() = %q", r.Canonical.Message())
	}
}

func TestRemoteErrorStatusCode(t *testing.T) {
	r := newRemoteErrorFixture()
	if r.StatusCode() != 503 {
		t.Errorf("StatusCode() = %d, want 503", r.StatusCode())
	}
}

// RemoteError is a leaf error: it has no cause and Unwrap returns nil.
// The Canonical view is intentionally not on the errors.As chain.
func TestRemoteErrorIsLeaf(t *testing.T) {
	r := newRemoteErrorFixture()
	if errors.Unwrap(r) != nil {
		t.Errorf("errors.Unwrap(*RemoteError) = %v, want nil", errors.Unwrap(r))
	}
}

// Canonical is a parallel view, not a chain step: errors.As to *AppError
// must NOT succeed via a RemoteError. Consumers access the canonical view
// explicitly through r.Canonical.
func TestRemoteErrorCanonicalNotInChain(t *testing.T) {
	r := newRemoteErrorFixture()
	var ae *AppError
	if errors.As(error(r), &ae) {
		t.Error("errors.As(remoteErr, &*AppError) should NOT succeed; " +
			"Canonical is a parallel view, not in the cause chain")
	}
}

func TestRemoteErrorErrorsAsRecoversRemoteError(t *testing.T) {
	r := newRemoteErrorFixture()
	wrapped := fmt.Errorf("calling user-service: %w", r) // pretend a caller wraps

	var re *RemoteError
	if !errors.As(wrapped, &re) {
		t.Fatal("errors.As should find *RemoteError in chain")
	}
	if re != r {
		t.Error("recovered RemoteError is not the original instance")
	}
}

func TestRemoteErrorStringFormat(t *testing.T) {
	r := newRemoteErrorFixture()
	got := r.Error()

	if !strings.HasPrefix(got, "RemoteError(") {
		t.Errorf("Error() should start with RemoteError(, got %q", got)
	}

	for _, want := range []string{
		"service=user-service",
		"operation=GetUser",
		"status=503",
		"bodyCode=DEGRADED",
		"retryAfter=30s",
		"code=UNAVAILABLE(14)",
		"message='user-service call failed'",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Error() missing %q\ngot: %s", want, got)
		}
	}
}

func TestRemoteErrorEvent(t *testing.T) {
	r := newRemoteErrorFixture()
	if got, want := r.Event(), "user-service.GetUser"; got != want {
		t.Errorf("Event() = %q, want %q", got, want)
	}
}

// RemoteError.Event() is independent of the Canonical's Event field:
// whatever event was passed to the Canonical's factory (event is now a
// required positional arg, so the caller MUST pass something), r.Event()
// still derives from Service.Operation. Convention: callers commonly pass
// Service.Operation to the Canonical so r.Canonical.Event() reads the
// same as r.Event(), but the value on the Canonical is never consulted
// by r.Event().
func TestRemoteErrorEventIndependentOfCanonical(t *testing.T) {
	r := &RemoteError{
		Canonical: NewUnavailable("unrelated.event", WithMessage("call failed")),
		Service:   "user-service",
		Operation: "GetUser",
		Response:  &Response{StatusCode: 503},
	}
	if got, want := r.Event(), "user-service.GetUser"; got != want {
		t.Errorf("Event() = %q, want %q", got, want)
	}
	if r.Canonical.Event() != "unrelated.event" {
		t.Errorf("Canonical.Event() not preserved (got %q)", r.Canonical.Event())
	}
}

func TestRemoteErrorAddNoteOnCanonical(t *testing.T) {
	r := newRemoteErrorFixture()
	r.Canonical.AddNote("during checkout")
	if !strings.HasPrefix(r.Canonical.Message(), "during checkout -> ") {
		t.Errorf("AddNote should prepend context; got %q", r.Canonical.Message())
	}
}
