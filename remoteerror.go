package apperror

import (
	"fmt"
	"time"
)

// RemoteError represents a failure observed when this application called
// another service over the network and **the server returned a response**.
// Any status code (2xx/4xx/5xx) or any in-band error envelope in the body
// qualifies, as long as a response was received. It covers both
// intra-company services and third-party APIs.
//
// # What is NOT a RemoteError
//
// Transport-layer failures — DNS resolution failure, connection refused or
// reset, TLS handshake failure, request timing out before any bytes were
// received — are NOT RemoteErrors. The remote service never "responded",
// so there is nothing for RemoteError to model. Represent those as plain
// AppError (typically NewUnavailable) wrapping the transport error via
// WithCause. The two error kinds are distinguished by *type*, not by
// fields inside one shared type.
//
// # Relationship to AppError
//
// RemoteError is NOT a subtype of AppError. The Canonical field is NOT a
// cause. It is a parallel, normalized view of the same RemoteError —
// the client at the call boundary picks one of this package's Code
// constants (e.g. CodeNotFound, CodeTooManyRequests) so cross-cutting
// logic has a stable taxonomy to branch on. Access it explicitly via
// r.Canonical; do not expect errors.As(err, &appErr) to recover it.
//
// # Three views of one error
//
// A RemoteError exposes three layers of error information, all describing
// the same failure from different angles:
//
//   - Canonical (Canonical.Code(), etc.) — our normalized taxonomy. Use
//     for retry / circuit breaker / log aggregation keys.
//
//   - Protocol (StatusCode, surfaced from Response.StatusCode) — the
//     transport-protocol status (HTTP status / gRPC status). Use for
//     retry decisions that depend on transport class.
//
//   - Remote application (BodyCode, BodyMessage) — the foreign
//     service's own application-level signals parsed from the response
//     body. Preserved for forensics and ops/runbook references; not the
//     right thing to branch on for cross-cutting logic.
//
// Conventions
//
//   - Canonical must be non-nil.
//   - Response must be non-nil — that's the precondition that makes this a
//     RemoteError in the first place.
//   - Do NOT call WithCause on the Canonical. RemoteError has no cause.
//   - The event you pass to the Canonical's factory is preserved on
//     r.Canonical.Event() but is NEVER consulted by r.Event() (which
//     always derives from Service.Operation). Convention: pass
//     Service+"."+Operation so r.Canonical.Event() and r.Event() agree
//     when read in isolation.
//   - Do NOT set the Canonical's Details. Remote-side structured info
//     lives on this struct's typed fields; raw payload lives in
//     Response.Body.
//   - Operation is a logical operation name (e.g. "GetUser"), not an HTTP
//     method + path.
type RemoteError struct {
	// Canonical is the normalized application-level view of this failure.
	// See package-level Code constants. Must be non-nil.
	Canonical *AppError

	// Service is the logical name of the remote service called
	// (e.g. "user-service", "stripe").
	Service string

	// Operation is the logical name of the operation invoked
	// (e.g. "GetUser", "CreateCharge"). Not an HTTP method + path.
	Operation string

	// Request is the captured outbound request. Often nil to avoid logging
	// sensitive bodies.
	Request *Request

	// Response is the captured response. Must be non-nil — its presence is
	// the precondition for using RemoteError.
	Response *Response

	// BodyCode is the remote service's application-level error code, parsed
	// from Response.Body (e.g. "card_declined"). Empty when the remote
	// didn't supply one. Useful as a low-cardinality aggregation key for
	// observability.
	BodyCode string

	// BodyMessage is the remote service's original error message, parsed
	// from Response.Body. Preserved for forensics.
	BodyMessage string

	// RetryAfter is the duration the remote service asks us to wait before
	// retrying, normalized from whatever form the remote used (HTTP
	// Retry-After header, a body field, etc.). Zero means the remote did
	// not provide a hint.
	RetryAfter time.Duration
}

// StatusCode returns the protocol-level status code of the captured
// Response. Response is required to be non-nil; if it isn't, this will
// panic (which is the correct behavior — it surfaces a convention violation
// loudly).
func (r *RemoteError) StatusCode() int {
	return r.Response.StatusCode
}

// Event returns "Service.Operation", suitable as a single low-cardinality
// key in structured logs or as a span/event name. For dashboards that need
// to slice by service or operation independently, log Service and Operation
// as separate fields too.
func (r *RemoteError) Event() string {
	return r.Service + "." + r.Operation
}

// RemoteError is a leaf error: it has no Unwrap because it has no cause.
// The Canonical view is NOT in any errors.Is / errors.As chain; reach it
// directly via r.Canonical.

// Error implements the error interface.
func (r *RemoteError) Error() string { return r.String() }

// String formats the error for human consumption. Output is RemoteError-
// shaped (not AppError-shaped) so the outer fields are visible.
func (r *RemoteError) String() string {
	var code Code
	var msg string
	if r.Canonical != nil {
		code = r.Canonical.Code()
		msg = r.Canonical.Message()
	}
	return fmt.Sprintf(
		"RemoteError(service=%s, operation=%s, status=%d, bodyCode=%s, retryAfter=%s, code=%s, message='%s')",
		r.Service, r.Operation, r.StatusCode(), r.BodyCode, r.RetryAfter, code, msg,
	)
}
