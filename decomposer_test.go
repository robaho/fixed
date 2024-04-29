//go:build !sql_scanner
// +build !sql_scanner

package fixed

import (
	"testing"
)

func TestDecomposerRoundTrip(t *testing.T) {
	list := []struct {
		N string // Name.
		S string // String value.
		E bool   // Expect an error.
	}{
		{N: "Zero", S: "0"},
		{N: "Normal-1", S: "123.456"},
		{N: "Normal-2", S: "-123.456"},
	}
	for _, item := range list {
		t.Run(item.N, func(t *testing.T) {
			d := NewS(item.S)
			if d.IsNaN() {
				t.Fatal("failed to parse number")
			}
			set := Fixed(0)
			err := set.Compose(d.Decompose(nil))
			if err == nil && item.E {
				t.Fatal("expected error, got <nil>")
			}
			if err != nil && !item.E {
				t.Fatalf("unexpected error: %v", err)
			}
			if set.Cmp(d) != 0 {
				t.Fatalf("values incorrect, got %v want %v (%s)", set, d, item.S)
			}
		})
	}
}

func TestDecomposerCompose(t *testing.T) {
	list := []struct {
		N string // Name.
		S string // String value.

		Form byte // Form
		Neg  bool
		Coef []byte // Coefficent
		Exp  int32

		Err bool // Expect an error.
	}{
		{N: "Zero", S: "0", Coef: nil, Exp: 0},
		{N: "Normal-1", S: "123.456", Coef: []byte{0x01, 0xE2, 0x40}, Exp: -3},
		{N: "Neg-1", S: "-123.456", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: -3},
		{N: "PosExp-1", S: "123456000", Coef: []byte{0x01, 0xE2, 0x40}, Exp: 3},
		{N: "PosExp-2", S: "-123456000", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: 3},
		{N: "AllDec-1", S: "0.123456", Coef: []byte{0x01, 0xE2, 0x40}, Exp: -6},
		{N: "AllDec-2", S: "-0.123456", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: -6},
		{N: "TooSmall-1", S: "-0.00123456", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: -8, Err: true},
		{N: "LeadingZero-1", S: "123.456", Coef: []byte{0, 0, 0, 0, 0, 0, 0, 0x01, 0xE2, 0x40}, Exp: -3},
		{N: "NaN-1", S: "NaN", Form: 2},
	}

	for _, item := range list {
		t.Run(item.N, func(t *testing.T) {
			d := Fixed(0)
			err := d.Compose(item.Form, item.Neg, item.Coef, item.Exp)
			if err != nil && !item.Err {
				t.Fatalf("unexpected error, got %v", err)
			}
			if item.Err {
				if err == nil {
					t.Fatal("expected error, got <nil>")
				}
				return
			}
			if s := d.String(); s != item.S {
				t.Fatalf("unexpected value, got %q want %q", s, item.S)
			}
		})
	}
}
