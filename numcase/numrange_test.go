package numcase

import (
	"strings"
	"testing"
)

func TestNumRangeInitValid(t *testing.T) {
	r1, err := New(0, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r1.Start() != 0 || r1.End() != 5 {
		t.Errorf("got [%d, %d], want [0, 5]", r1.Start(), r1.End())
	}

	r2, err := New(10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r2.Start() != 10 || r2.End() != 10 {
		t.Errorf("got [%d, %d], want [10, 10]", r2.Start(), r2.End())
	}
}

func TestNumRangeInitInvalid(t *testing.T) {
	cases := []struct {
		start, end int
		match      string
	}{
		{-1, 5, "start < 0"},
		{0, -1, "end < 0"},
		{10, 5, "end < start"},
	}
	for _, c := range cases {
		_, err := New(c.start, c.end)
		if err == nil {
			t.Errorf("New(%d, %d) should error", c.start, c.end)
			continue
		}
		if !strings.Contains(err.Error(), c.match) {
			t.Errorf("New(%d, %d) err = %q, want substring %q",
				c.start, c.end, err.Error(), c.match)
		}
	}
}

func TestNumRangeInclude(t *testing.T) {
	r := MustNew(0, 5)
	cases := []struct {
		num  int
		want bool
	}{
		{0, true}, {3, true}, {5, true}, {-1, false}, {6, false},
	}
	for _, c := range cases {
		if got := r.Include(c.num); got != c.want {
			t.Errorf("Include(%d) = %v, want %v", c.num, got, c.want)
		}
	}
}

func TestNumRangeIncludeRange(t *testing.T) {
	r1 := MustNew(0, 10)
	r2 := MustNew(2, 8)
	r3 := MustNew(5, 15)
	if !r1.IncludeRange(r2) {
		t.Error("r1 should include r2")
	}
	if r1.IncludeRange(r3) {
		t.Error("r1 should not include r3")
	}
	if !r1.IncludeRange(r1) {
		t.Error("r1 should include itself")
	}
}

func TestNumRangeOverlap(t *testing.T) {
	if !MustNew(0, 5).Overlap(MustNew(5, 10)) {
		t.Error("[0,5] and [5,10] should overlap")
	}
	if MustNew(0, 5).Overlap(MustNew(6, 10)) {
		t.Error("[0,5] and [6,10] should not overlap")
	}
}

func TestNumRangeString(t *testing.T) {
	if got := MustNew(0, 5).String(); got != "[0, 5]" {
		t.Errorf("String() = %q, want %q", got, "[0, 5]")
	}
}

func TestNumRangeGoString(t *testing.T) {
	if got := MustNew(0, 5).GoString(); got != "NumRange(start=0, end=5)" {
		t.Errorf("GoString() = %q, want %q", got, "NumRange(start=0, end=5)")
	}
}
