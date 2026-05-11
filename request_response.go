package apperror

// Request captures an outbound HTTP/RPC request associated with an error.
// Capture is optional; clients often omit it (or omit just the Body) to
// avoid leaking sensitive request payloads through logs.
type Request struct {
	URL     string
	Method  string
	Headers map[string][]string
	Body    []byte // raw bytes
}

// Response captures an HTTP/RPC response associated with an error. Its
// presence is what makes a failure a RemoteError (the server received the
// request and replied); transport-layer failures where no response was
// received are NOT modeled as RemoteError — represent those as plain
// AppError, see the RemoteError doc.
type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte // raw bytes; parsed views live on RemoteError fields
}
