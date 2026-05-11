package apperror

import (
	"errors"
	"strings"
	"testing"
)

// errorCase mirrors the Python ErrorCase used in tests: a Case identified by
// a string, with an associated op-status code that is not asserted on.
type errorCase struct {
	id     string
	opCode Code
}

func (e errorCase) Identifier() string { return e.id }

// connection raises a low-level RuntimeError-like error.
type connection struct{}

func (connection) ExecSQL(sql string) error {
	return errors.New("Network error")
}

type dbClient struct{ conn connection }

func (d dbClient) Insert(data string) error {
	if err := d.conn.ExecSQL("insert into ..."); err != nil {
		return NewInternalError("DB insert failed", WithCause(err))
	}
	return nil
}

type service struct{ db dbClient }

func (s service) Save(data string) error { return s.db.Insert(data) }

type api struct{ s service }

func (a api) Create(data string) error { return a.s.Save(data) }

func TestCheckAppErrorCause(t *testing.T) {
	a := api{s: service{db: dbClient{conn: connection{}}}}
	err := a.Create("test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *AppError, got %T", err)
	}
	cause := appErr.Cause()
	if cause == nil {
		t.Fatal("expected cause to be set")
	}
	if cause.Error() != "Network error" {
		t.Errorf("cause.Error() = %q, want %q", cause.Error(), "Network error")
	}

	// Unwrap path also works for errors.Is.
	if !errors.Is(err, cause) {
		t.Error("errors.Is(err, cause) should be true")
	}
}

func TestCreateAppErrorWithCause(t *testing.T) {
	low := errors.New("low-level error in func_a")
	appErr := NewInternalError("higher-level error msg", WithCause(low))
	appErr.AddErrCtx("additional err context...")
	if !errors.Is(appErr, low) {
		t.Errorf("errors.Is(appErr, low) = false, want true")
	}
}

func TestAddErrCtx(t *testing.T) {
	e := NewInternalError("initial error message")
	e.AddErrCtx("first context")
	e.AddErrCtx("second context")
	want := "second context -> first context -> initial error message"
	if e.Message() != want {
		t.Errorf("Message() = %q, want %q", e.Message(), want)
	}
	if !strings.Contains(e.Error(), want) {
		t.Errorf("Error() should contain %q, got %q", want, e.Error())
	}
}

func TestAppErrorString(t *testing.T) {
	cases := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "basic",
			err:  NewInternalError("internal error"),
			want: "AppError(code=INTERNAL_ERROR(13), case=None, event=None, message='internal error', details=None)",
		},
		{
			name: "with case",
			err:  NewInternalError("internal error", WithCase(errorCase{id: "1001", opCode: CodeInternalError})),
			want: "AppError(code=INTERNAL_ERROR(13), case=1001, event=None, message='internal error', details=None)",
		},
		{
			name: "with details",
			err:  NewInternalError("internal error", WithDetails(map[string]string{"key": "value"})),
			want: "AppError(code=INTERNAL_ERROR(13), case=None, event=None, message='internal error', details=map[key:value])",
		},
		{
			name: "with event",
			err:  NewInternalError("internal error", WithEvent("user.signup")),
			want: "AppError(code=INTERNAL_ERROR(13), case=None, event=user.signup, message='internal error', details=None)",
		},
		{
			name: "full",
			err: NewInternalError("internal error",
				WithCase(errorCase{id: "1001", opCode: CodeInternalError}),
				WithDetails(map[string]string{"key": "value"}),
				WithEvent("user.signup")),
			want: "AppError(code=INTERNAL_ERROR(13), case=1001, event=user.signup, message='internal error', details=map[key:value])",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.err.String(); got != c.want {
				t.Errorf("String() = %q\nwant      %q", got, c.want)
			}
		})
	}
}

func TestBuildErrorWithOptions(t *testing.T) {
	e := NewInternalError("internal error")
	if e.Code() != CodeInternalError || e.Message() != "internal error" || e.Details() != nil {
		t.Errorf("unexpected basic error: %v", e)
	}

	d := map[string]any{"key1": "value1"}
	e = NewInternalError("internal error", WithDetails(d))
	got, ok := e.Details().(map[string]any)
	if !ok || got["key1"] != "value1" {
		t.Errorf("unexpected details: %v", e.Details())
	}

	d2 := map[string]any{"key1": "value1", "extra_info": map[string]string{"key2": "value2"}}
	e = NewInternalError("internal error", WithDetails(d2))
	got, ok = e.Details().(map[string]any)
	if !ok || got["key1"] != "value1" {
		t.Errorf("unexpected nested details: %v", e.Details())
	}
}

func TestNewIllegalArg(t *testing.T) {
	e := NewIllegalArg("illegal arg")
	if e.Code() != CodeIllegalArg || e.Message() != "illegal arg" {
		t.Errorf("unexpected illegal-arg error: %v", e)
	}

	e = NewIllegalArg("Test error")
	want := "AppError(code=ILLEGAL_ARG(29), case=None, event=None, message='Test error', details=None)"
	if got := e.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestAddMoreErrCtx(t *testing.T) {
	validate := func() error { return NewIllegalArg("illegal arg") }

	err := validate()
	var appErr *AppError
	if errors.As(err, &appErr) {
		appErr.AddErrCtx("Error while executing ...")
	}
	want := "AppError(code=ILLEGAL_ARG(29), case=None, event=None, message='Error while executing ... -> illegal arg', details=None)"
	if got := appErr.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestFactoryFallsBackToCodeDescriptionForEmptyMessage(t *testing.T) {
	e := NewInternalError("")
	if e.Message() != CodeInternalError.Description() {
		t.Errorf("empty message should fall back to description, got %q", e.Message())
	}
}

func TestAppErrorEvent(t *testing.T) {
	// Default: empty.
	e := NewInternalError("boom")
	if e.Event() != "" {
		t.Errorf("default Event() = %q, want empty", e.Event())
	}

	// WithEvent sets it.
	e = NewInternalError("boom", WithEvent("user.signup"))
	if e.Event() != "user.signup" {
		t.Errorf("Event() = %q, want %q", e.Event(), "user.signup")
	}
}
