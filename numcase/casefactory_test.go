package numcase

import (
	"testing"
)

func newDefaultFactory(t *testing.T) *CaseFactory {
	t.Helper()
	f, err := NewCaseFactory(FactoryConfig{
		NumDigitsForAppCode:    2,
		NumDigitsForModuleCode: 2,
		NumDigitsForCaseCode:   3,
		CodeMapper:             NewDefaultCodeMapper(),
		AppCode:                1,
		ModuleCode:             2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return f
}

func TestCaseFactoryInit(t *testing.T) {
	f := newDefaultFactory(t)
	if f.AppCode() != 1 || f.ModuleCode() != 2 {
		t.Errorf("AppCode/ModuleCode = (%d, %d), want (1, 2)", f.AppCode(), f.ModuleCode())
	}
	if f.NumDigitsForAppCode() != 2 || f.NumDigitsForModuleCode() != 2 || f.NumDigitsForCaseCode() != 3 {
		t.Errorf("digit widths mismatch")
	}

	bad := []FactoryConfig{
		{NumDigitsForAppCode: -1, NumDigitsForModuleCode: 2, NumDigitsForCaseCode: 3, CodeMapper: NewDefaultCodeMapper()},
		{NumDigitsForAppCode: 2, NumDigitsForModuleCode: -1, NumDigitsForCaseCode: 3, CodeMapper: NewDefaultCodeMapper()},
		{NumDigitsForAppCode: 2, NumDigitsForModuleCode: 2, NumDigitsForCaseCode: -1, CodeMapper: NewDefaultCodeMapper()},
		{NumDigitsForAppCode: 2, NumDigitsForModuleCode: 2, NumDigitsForCaseCode: 3, CodeMapper: NewDefaultCodeMapper(), AppCode: -1},
		{NumDigitsForAppCode: 2, NumDigitsForModuleCode: 2, NumDigitsForCaseCode: 3, CodeMapper: NewDefaultCodeMapper(), ModuleCode: -1},
	}
	for i, cfg := range bad {
		if _, err := NewCaseFactory(cfg); err == nil {
			t.Errorf("case %d: expected error, got nil", i)
		}
	}
}

func TestBuildCaseID(t *testing.T) {
	f := newDefaultFactory(t)
	if got := f.BuildCaseID(123); got != "01_02_123" {
		t.Errorf("BuildCaseID(123) = %q, want %q", got, "01_02_123")
	}
	if got := f.BuildCaseID(1); got != "01_02_001" {
		t.Errorf("BuildCaseID(1) = %q, want %q", got, "01_02_001")
	}

	f2, err := NewCaseFactory(FactoryConfig{
		NumDigitsForAppCode:    0,
		NumDigitsForModuleCode: 0,
		NumDigitsForCaseCode:   3,
		CodeMapper:             NewDefaultCodeMapper(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := f2.BuildCaseID(123); got != "123" {
		t.Errorf("BuildCaseID(123) = %q, want %q", got, "123")
	}
}

func TestCreateCases(t *testing.T) {
	f := newDefaultFactory(t)

	c, err := f.NewIllegalInput(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.AppCode() != 1 || c.ModuleCode() != 2 || c.CaseCode() != 1 {
		t.Errorf("got (%d, %d, %d), want (1, 2, 1)", c.AppCode(), c.ModuleCode(), c.CaseCode())
	}

	c, err = f.NewNotFound(101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.CaseCode() != 101 {
		t.Errorf("CaseCode() = %d, want 101", c.CaseCode())
	}

	if _, err := f.NewIllegalInput(51); err == nil {
		t.Error("expected error for case_code outside IllegalInput segment")
	}
	if _, err := f.NewNotFound(1); err == nil {
		t.Error("expected error for case_code outside NotFound segment")
	}
}

func TestPadLeftZeros(t *testing.T) {
	cases := []struct {
		num, minLen int
		want        string
	}{
		{1, 3, "001"},
		{12, 3, "012"},
		{123, 3, "123"},
		{1234, 3, "1234"},
	}
	for _, c := range cases {
		if got := padLeftZeros(c.num, c.minLen); got != c.want {
			t.Errorf("padLeftZeros(%d, %d) = %q, want %q", c.num, c.minLen, got, c.want)
		}
	}
}
