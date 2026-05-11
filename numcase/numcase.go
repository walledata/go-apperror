package numcase

// NumCase is a specific error condition represented by a numerical code.
//
// The identifier format is: {app_code}_{module_code}_{case_code}, for
// example: "1_1_1000". When AppCode or ModuleCode parts are configured with
// zero digits, those parts (and their separators) are omitted.
type NumCase struct {
	appCode    int
	moduleCode int
	caseCode   int
	identifier string
}

// AppCode returns the application code.
func (c NumCase) AppCode() int { return c.appCode }

// ModuleCode returns the module code.
func (c NumCase) ModuleCode() int { return c.moduleCode }

// CaseCode returns the per-case numeric code.
func (c NumCase) CaseCode() int { return c.caseCode }

// Identifier returns the assembled identifier string. NumCase satisfies the
// apperror.Case interface via this method.
func (c NumCase) Identifier() string { return c.identifier }

// String returns the identifier (used by fmt for default printing).
func (c NumCase) String() string { return c.identifier }

// Equal reports whether c and other have the same identifier.
func (c NumCase) Equal(other NumCase) bool {
	return c.identifier == other.identifier
}
