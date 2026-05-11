package apperror

var codeToHTTPStatus = map[Code]HTTPStatus{
	CodeOK:                   StatusOK,
	CodeIllegalInput:         StatusBadRequest,
	CodeFailedPrecondition:   StatusBadRequest,
	CodeOutOfRange:           StatusBadRequest,
	CodeUnauthenticated:      StatusUnauthorized,
	CodePermissionDenied:     StatusForbidden,
	CodeNotFound:             StatusNotFound,
	CodeOpConflict:           StatusConflict,
	CodeAlreadyExists:        StatusConflict,
	CodeTooManyRequests:      StatusTooManyRequests,
	CodeOpCancelled:          StatusClientClosedRequest,
	CodeIllegalState:         StatusInternalServerError,
	CodeUnknownError:         StatusInternalServerError,
	CodeInternalError:        StatusInternalServerError,
	CodeUnimplemented:        StatusNotImplemented,
	CodeUnavailable:          StatusServiceUnavailable,
	CodeTimeout:              StatusTimeout,
	CodeAuthorizationExpired: StatusUnauthorized,
	CodeIllegalArg:           StatusInternalServerError,
}

// HTTPStatusFor returns the HTTPStatus that the given Code maps to.
// The second return value is false when no mapping is defined.
func HTTPStatusFor(code Code) (HTTPStatus, bool) {
	s, ok := codeToHTTPStatus[code]
	return s, ok
}

var httpStatusToCode = map[HTTPStatus]Code{
	StatusOK:                  CodeOK,
	StatusBadRequest:          CodeIllegalInput,
	StatusUnauthorized:        CodeUnauthenticated,
	StatusForbidden:           CodePermissionDenied,
	StatusNotFound:            CodeNotFound,
	StatusConflict:            CodeAlreadyExists,
	StatusTooManyRequests:     CodeTooManyRequests,
	StatusClientClosedRequest: CodeOpCancelled,
	StatusInternalServerError: CodeInternalError,
	StatusNotImplemented:      CodeUnimplemented,
	StatusServiceUnavailable:  CodeUnavailable,
	StatusTimeout:             CodeTimeout,
}

// OpCodeFor returns the operation status Code that the given HTTP status code
// maps to. Unknown statuses map to CodeUnknownError.
func OpCodeFor(httpStatus int) Code {
	if c, ok := httpStatusToCode[HTTPStatus(httpStatus)]; ok {
		return c
	}
	return CodeUnknownError
}
