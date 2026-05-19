package apperror

// Case represents a specific error condition,
// for example: purchase_limit_exceeded, insufficient_inventory.
//
// Most errors don't need a Case — Code already conveys the failure
// category. Define one only when callers must branch on a specific
// business condition within that category (e.g. distinguishing
// "email already taken" from "phone already taken" under AlreadyExists
// so the UI can suggest password recovery).
type Case interface {
	// Identifier returns a string that uniquely identifies this error case.
	// It can be a numerical value or a descriptive title/name. For example,
	// numerical values: "1000", "1_1_1000"; or descriptive names:
	// "purchase_limit_exceeded".
	Identifier() string
}

// StrCase is a Case identified by some words or a phrase, for example
// "purchase_limit_exceeded".
type StrCase struct {
	id string
}

// NewStrCase creates a StrCase with the given identifier.
func NewStrCase(id string) StrCase {
	return StrCase{id: id}
}

// Identifier returns the identifier of the case.
func (c StrCase) Identifier() string {
	return c.id
}

// String returns the case identifier (used by fmt for default printing).
func (c StrCase) String() string {
	return c.id
}
