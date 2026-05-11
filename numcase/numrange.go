// Package numcase implements numeric Case identifiers, their valid ranges,
// and a factory for assembling them under an application/module/case scheme.
package numcase

import (
	"errors"
	"fmt"
)

// NumRange represents a number range from Start to End, both inclusive.
type NumRange struct {
	start int
	end   int
}

// New constructs a NumRange. It returns an error if start < 0, end < 0, or
// end < start.
func New(start, end int) (NumRange, error) {
	if start < 0 {
		return NumRange{}, errors.New("start < 0")
	}
	if end < 0 {
		return NumRange{}, errors.New("end < 0")
	}
	if end < start {
		return NumRange{}, errors.New("end < start")
	}
	return NumRange{start: start, end: end}, nil
}

// MustNew constructs a NumRange and panics on invalid input. Useful for
// package-level static configuration where invalid values are programmer
// errors.
func MustNew(start, end int) NumRange {
	r, err := New(start, end)
	if err != nil {
		panic(err)
	}
	return r
}

// Start returns the inclusive lower bound.
func (r NumRange) Start() int { return r.start }

// End returns the inclusive upper bound.
func (r NumRange) End() int { return r.end }

// Include reports whether num is within this range (inclusive).
func (r NumRange) Include(num int) bool {
	return r.start <= num && num <= r.end
}

// IncludeRange reports whether other is fully contained within r.
func (r NumRange) IncludeRange(other NumRange) bool {
	return r.start <= other.start && other.end <= r.end
}

// Overlap reports whether r and other share any value.
func (r NumRange) Overlap(other NumRange) bool {
	return r.end >= other.start && other.end >= r.start
}

// Equal reports value-equality with other.
func (r NumRange) Equal(other NumRange) bool {
	return r.start == other.start && r.end == other.end
}

// String returns "[start, end]".
func (r NumRange) String() string {
	return fmt.Sprintf("[%d, %d]", r.start, r.end)
}

// GoString returns a debug-friendly representation.
func (r NumRange) GoString() string {
	return fmt.Sprintf("NumRange(start=%d, end=%d)", r.start, r.end)
}
