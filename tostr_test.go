package fixed_test

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	. "github.com/PKartaviy/fixed"
)

// TestTostrViaString thoroughly tests the tostr() implementation via String()
func TestTostrViaString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		// Zero
		{"zero", "0", "0"},
		{"zero from 0.0", "0.0", "0"},
		{"zero from 0.000", "0.000", "0"},

		// Simple positive integers
		{"one", "1", "1"},
		{"nine", "9", "9"},
		{"ten", "10", "10"},
		{"hundred", "100", "100"},
		{"thousand", "1000", "1000"},
		{"million", "1000000", "1000000"},
		{"billion", "1000000000", "1000000000"},

		// Simple negative integers
		{"neg one", "-1", "-1"},
		{"neg nine", "-9", "-9"},
		{"neg ten", "-10", "-10"},
		{"neg hundred", "-100", "-100"},
		{"neg thousand", "-1000", "-1000"},
		{"neg million", "-1000000", "-1000000"},

		// Simple decimals
		{"0.1", "0.1", "0.1"},
		{"0.5", "0.5", "0.5"},
		{"0.9", "0.9", "0.9"},
		{"0.01", "0.01", "0.01"},
		{"0.001", "0.001", "0.001"},
		{"0.0001", "0.0001", "0.0001"},
		{"0.12345", "0.12345", "0.12345"},

		// Negative decimals
		{"-0.1", "-0.1", "-0.1"},
		{"-0.5", "-0.5", "-0.5"},
		{"-0.01", "-0.01", "-0.01"},
		{"-0.001", "-0.001", "-0.001"},
		{"-0.12345", "-0.12345", "-0.12345"},

		// Mixed integer and decimal
		{"1.1", "1.1", "1.1"},
		{"1.5", "1.5", "1.5"},
		{"12.34", "12.34", "12.34"},
		{"123.456", "123.456", "123.456"},
		{"1234.5678", "1234.5678", "1234.5678"},
		{"123456789.123456789", "123456789.123456789", "123456789.123456789"},
		{"-1.1", "-1.1", "-1.1"},
		{"-12.34", "-12.34", "-12.34"},
		{"-123.456", "-123.456", "-123.456"},

		// Trailing zeros should be stripped by String()
		{"trailing zeros", "1.10", "1.1"},
		{"trailing zeros 2", "1.100", "1.1"},
		{"trailing zeros 3", "100.00", "100"},
		{"trailing zeros 4", "0.10", "0.1"},
		{"trailing zeros 5", "0.010", "0.01"},

		// Leading zeros in fractional part (critical for zero-padding)
		{"leading zero frac", "1.01", "1.01"},
		{"leading zero frac 2", "1.001", "1.001"},
		{"leading zero frac 3", "1.0001", "1.0001"},
		{"leading zero frac 4", "1.00001", "1.00001"},
		{"leading zero frac 5", "1.000001", "1.000001"},
		{"leading zero frac 6", "1.0000001", "1.0000001"},
		{"leading zero frac 7", "1.00000001", "1.00000001"},
		{"leading zero frac 8", "1.000000001", "1.000000001"},
		{"leading zero frac 9", "1.0000000001", "1.0000000001"},
		{"leading zero frac 17", "1.00000000000000001", "1.00000000000000001"},

		// Minimum fractional value (1 in the lowest decimal place)
		{"min frac", "0.000000000000000001", "0.000000000000000001"},
		{"-min frac", "-0.000000000000000001", "-0.000000000000000001"},
		{"1 + min frac", "1.000000000000000001", "1.000000000000000001"},

		// Maximum values (18 integer digits)
		{"max hi", "999999999999999999", "999999999999999999"},
		{"-max hi", "-999999999999999999", "-999999999999999999"},

		// Max fractional
		{"max frac", "0.999999999999999999", "0.999999999999999999"},
		{"-max frac", "-0.999999999999999999", "-0.999999999999999999"},

		// Max combined
		{"max combined", "999999999999999999.999999999999999999", "999999999999999999.999999999999999999"},
		{"-max combined", "-999999999999999999.999999999999999999", "-999999999999999999.999999999999999999"},

		// Powers of 10
		{"10^1", "10", "10"},
		{"10^2", "100", "100"},
		{"10^3", "1000", "1000"},
		{"10^6", "1000000", "1000000"},
		{"10^9", "1000000000", "1000000000"},
		{"10^12", "1000000000000", "1000000000000"},
		{"10^15", "1000000000000000", "1000000000000000"},
		{"10^17", "100000000000000000", "100000000000000000"},

		// Single digits in various positions
		{"hi=1 lo=0", "1", "1"},
		{"hi=0 lo=1", "0.000000000000000001", "0.000000000000000001"},
		{"hi=1 lo=1", "1.000000000000000001", "1.000000000000000001"},
		{"hi=9 lo=9", "9.000000000000000009", "9.000000000000000009"},

		// Patterns that stress zero-padding
		{"100000000000000000.1", "100000000000000000.1", "100000000000000000.1"},
		{"1.100000000000000000", "1.1", "1.1"},

		// NaN
		{"NaN", "NaN", "NaN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewS(tt.input)
			got := f.String()
			if got != tt.expect {
				t.Errorf("NewS(%q).String() = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

// TestTostrViaStringN tests StringN() output for various decimal place counts
func TestTostrViaStringN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		decimals int
		expect   string
	}{
		// Zero with various decimal places
		{"zero N=0", "0", 0, "0"},
		{"zero N=1", "0", 1, "0.0"},
		{"zero N=2", "0", 2, "0.00"},
		{"zero N=5", "0", 5, "0.00000"},
		{"zero N=18", "0", 18, "0.000000000000000000"},

		// Integer with decimal places
		{"1 N=0", "1", 0, "1"},
		{"1 N=1", "1", 1, "1.0"},
		{"1 N=2", "1", 2, "1.00"},
		{"1 N=5", "1", 5, "1.00000"},

		// Truncation (not rounding)
		{"1.129 N=2", "1.129", 2, "1.12"},
		{"1.999 N=2", "1.999", 2, "1.99"},
		{"1.999 N=1", "1.999", 1, "1.9"},
		{"1.999 N=0", "1.999", 0, "1"},

		// Exact decimal places
		{"1.12 N=2", "1.12", 2, "1.12"},
		{"1.123 N=3", "1.123", 3, "1.123"},

		// Padding with zeros
		{"1.1 N=5", "1.1", 5, "1.10000"},
		{"1.12 N=5", "1.12", 5, "1.12000"},

		// Negative values
		{"-1.129 N=2", "-1.129", 2, "-1.12"},
		{"-1.1 N=5", "-1.1", 5, "-1.10000"},
		{"-1 N=2", "-1", 2, "-1.00"},

		// Large values
		{"large N=2", "999999999999999999.999", 2, "999999999999999999.99"},
		{"large N=0", "999999999999999999.999", 0, "999999999999999999"},
		{"large N=18", "999999999999999999.999999999999999999", 18, "999999999999999999.999999999999999999"},

		// Min fractional
		{"min frac N=18", "0.000000000000000001", 18, "0.000000000000000001"},
		{"min frac N=17", "0.000000000000000001", 17, "0.00000000000000000"},
		{"min frac N=1", "0.000000000000000001", 1, "0.0"},

		// NaN
		{"NaN N=2", "NaN", 2, "NaN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewS(tt.input)
			got := f.StringN(tt.decimals)
			if got != tt.expect {
				t.Errorf("NewS(%q).StringN(%d) = %q, want %q", tt.input, tt.decimals, got, tt.expect)
			}
		})
	}
}

// TestTostrPointPosition verifies that the point position returned by tostr()
// is correct by checking that StringN works properly (it depends on point position)
func TestTostrPointPosition(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"zero", "0"},
		{"positive int", "123"},
		{"negative int", "-123"},
		{"positive decimal", "123.456"},
		{"negative decimal", "-123.456"},
		{"small decimal", "0.001"},
		{"negative small decimal", "-0.001"},
		{"large", "999999999999999999.999999999999999999"},
		{"negative large", "-999999999999999999.999999999999999999"},
		{"single digit", "1"},
		{"neg single digit", "-1"},
		{"min frac", "0.000000000000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewS(tt.input)
			// StringN(0) should give just the integer part — verifies point is correct
			sn0 := f.StringN(0)
			full := f.String()

			// The integer part from StringN(0) should be a prefix of String() output
			// (possibly followed by decimal point and more digits)
			if !strings.HasPrefix(full, sn0) && !strings.HasPrefix(full, sn0+".") {
				// For integers, String() drops the decimal, so they should be equal
				if full != sn0 {
					t.Errorf("StringN(0)=%q is not a prefix of String()=%q", sn0, full)
				}
			}

			// StringN(18) should give all 18 decimal places
			sn18 := f.StringN(18)
			if f.String() != "NaN" && len(sn18) < len(sn0)+19 { // +1 for dot + 18 digits
				t.Errorf("StringN(18)=%q too short (StringN(0)=%q)", sn18, sn0)
			}
		})
	}
}

// TestTostrFromNewF tests String() output for values constructed via NewF (float64)
func TestTostrFromNewF(t *testing.T) {
	tests := []struct {
		name   string
		input  float64
		expect string
	}{
		{"zero", 0.0, "0"},
		{"one", 1.0, "1"},
		{"neg one", -1.0, "-1"},
		{"0.5", 0.5, "0.5"},
		{"-0.5", -0.5, "-0.5"},
		{"0.1", 0.1, "0.1"},
		{"0.01", 0.01, "0.01"},
		{"0.001", 0.001, "0.001"},
		{"0.0001", 0.0001, "0.0001"},
		{"123.456", 123.456, "123.456"},
		{"-123.456", -123.456, "-123.456"},
		{"NaN", math.NaN(), "NaN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewF(tt.input)
			got := f.String()
			if got != tt.expect {
				t.Errorf("NewF(%v).String() = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

// TestTostrFromNewI tests String() for values constructed via NewI
func TestTostrFromNewI(t *testing.T) {
	tests := []struct {
		name     string
		val      int64
		decimals uint
		expect   string
	}{
		{"0,0", 0, 0, "0"},
		{"1,0", 1, 0, "1"},
		{"-1,0", -1, 0, "-1"},
		{"123,1", 123, 1, "12.3"},
		{"-123,1", -123, 1, "-12.3"},
		{"123,0", 123, 0, "123"},
		{"100,2", 100, 2, "1"},
		{"1,18", 1, 18, "0.000000000000000001"},
		{"-1,18", -1, 18, "-0.000000000000000001"},
		{"123456789012,9", 123456789012, 9, "123.456789012"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewI(tt.val, tt.decimals)
			got := f.String()
			if got != tt.expect {
				t.Errorf("NewI(%d, %d).String() = %q, want %q", tt.val, tt.decimals, got, tt.expect)
			}
		})
	}
}

// TestTostrMarshalJSON tests that MarshalJSON produces correct output
func TestTostrMarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"zero", "0", "0"},
		{"positive", "123.456", "123.456"},
		{"negative", "-123.456", "-123.456"},
		{"integer", "42", "42"},
		{"small decimal", "0.000000000000000001", "0.000000000000000001"},
		{"max value", "999999999999999999.999999999999999999", "999999999999999999.999999999999999999"},
		{"NaN", "NaN", "\"NaN\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewS(tt.input)
			got, err := f.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON error: %v", err)
			}
			if string(got) != tt.expect {
				t.Errorf("NewS(%q).MarshalJSON() = %q, want %q", tt.input, string(got), tt.expect)
			}
		})
	}
}

// TestTostrMarshalJSONRoundTrip tests JSON marshal/unmarshal round-trip
func TestTostrMarshalJSONRoundTrip(t *testing.T) {
	type wrapper struct {
		V Fixed `json:"v"`
	}

	values := []string{
		"0", "1", "-1",
		"0.1", "-0.1",
		"0.000000000000000001", "-0.000000000000000001",
		"123.456", "-123.456",
		"999999999999999999", "-999999999999999999",
		"999999999999999999.999999999999999999",
		"NaN",
	}

	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			original := NewS(v)
			w := wrapper{V: original}

			data, err := json.Marshal(w)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}

			var w2 wrapper
			err = json.Unmarshal(data, &w2)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if original.IsNaN() {
				if !w2.V.IsNaN() {
					t.Errorf("expected NaN after round-trip, got %v", w2.V)
				}
				return
			}

			if !w2.V.Equal(original) {
				t.Errorf("round-trip mismatch: %q -> marshal -> unmarshal -> %q", original.String(), w2.V.String())
			}
		})
	}
}

// TestTostrAfterArithmetic tests String() output after arithmetic operations
// to ensure tostr() handles all resulting internal states correctly
func TestTostrAfterArithmetic(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		f := NewS("0.1").Add(NewS("0.2"))
		got := f.String()
		if got != "0.3" {
			t.Errorf("0.1 + 0.2 = %q, want %q", got, "0.3")
		}
	})

	t.Run("subtraction to zero", func(t *testing.T) {
		f := NewS("1.5").Sub(NewS("1.5"))
		got := f.String()
		if got != "0" {
			t.Errorf("1.5 - 1.5 = %q, want %q", got, "0")
		}
	})

	t.Run("subtraction to negative", func(t *testing.T) {
		f := NewS("1").Sub(NewS("2"))
		got := f.String()
		if got != "-1" {
			t.Errorf("1 - 2 = %q, want %q", got, "-1")
		}
	})

	t.Run("multiplication", func(t *testing.T) {
		f := NewS("3").Mul(NewS("7"))
		got := f.String()
		if got != "21" {
			t.Errorf("3 * 7 = %q, want %q", got, "21")
		}
	})

	t.Run("division producing repeating decimal", func(t *testing.T) {
		f := NewS("1").Div(NewS("3"))
		got := f.String()
		if got != "0.333333333333333333" {
			t.Errorf("1 / 3 = %q, want %q", got, "0.333333333333333333")
		}
	})

	t.Run("division producing 2/3", func(t *testing.T) {
		f := NewS("2").Div(NewS("3"))
		got := f.String()
		if got != "0.666666666666666667" {
			t.Errorf("2 / 3 = %q, want %q", got, "0.666666666666666667")
		}
	})

	t.Run("large multiplication", func(t *testing.T) {
		f := NewS("999999999").Mul(NewS("999999999"))
		got := f.String()
		if got != "999999998000000001" {
			t.Errorf("999999999 * 999999999 = %q, want %q", got, "999999998000000001")
		}
	})
}

// TestTostrStringConsistency verifies String() and StringN(18) consistency
// (StringN(18) should be String() + trailing zeros if String() has fewer than 18 decimal places)
func TestTostrStringConsistency(t *testing.T) {
	values := []string{
		"0", "1", "-1",
		"0.1", "-0.1",
		"123.456", "-123.456",
		"0.000000000000000001",
		"999999999999999999.999999999999999999",
		"100", "1000000",
	}

	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			f := NewS(v)
			s := f.String()
			sn18 := f.StringN(18)

			if s == "NaN" {
				return
			}

			// Find the decimal point in both
			dotS := strings.Index(s, ".")
			dotSN := strings.Index(sn18, ".")

			if dotS == -1 {
				// String() has no decimal point — it's a pure integer
				// StringN(18) should be integer + "." + 18 zeros
				expected := s + "." + strings.Repeat("0", 18)
				if sn18 != expected {
					t.Errorf("StringN(18)=%q, want %q (String()=%q)", sn18, expected, s)
				}
			} else {
				// Both should share the same integer part
				if s[:dotS] != sn18[:dotSN] {
					t.Errorf("integer parts differ: String()=%q, StringN(18)=%q", s, sn18)
				}
				// StringN(18) fractional part should start with String()'s fractional part
				sFrac := s[dotS+1:]
				snFrac := sn18[dotSN+1:]
				if !strings.HasPrefix(snFrac, sFrac) {
					t.Errorf("fractional mismatch: String() frac=%q, StringN(18) frac=%q", sFrac, snFrac)
				}
				// Remaining chars in StringN(18) should all be zeros
				remainder := snFrac[len(sFrac):]
				if remainder != strings.Repeat("0", len(remainder)) {
					t.Errorf("StringN(18) has non-zero padding: %q", remainder)
				}
			}
		})
	}
}

// TestTostrNewSRoundTrip verifies that NewS(f.String()) == f for various values
func TestTostrNewSRoundTrip(t *testing.T) {
	values := []string{
		"0",
		"1", "-1",
		"0.1", "-0.1",
		"0.5", "-0.5",
		"0.000000000000000001", "-0.000000000000000001",
		"123456789.123456789", "-123456789.123456789",
		"999999999999999999", "-999999999999999999",
		"999999999999999999.999999999999999999", "-999999999999999999.999999999999999999",
		"100000000000000000", "100000000000000000.1",
		"1.000000000000000001",
	}

	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			f1 := NewS(v)
			s := f1.String()
			f2 := NewS(s)

			if !f1.Equal(f2) {
				t.Errorf("round-trip failed: %q -> String() -> %q -> NewS() -> %q", v, s, f2.String())
			}
		})
	}
}

// TestTostrStringLength ensures the output never exceeds expected bounds
func TestTostrStringLength(t *testing.T) {
	values := []string{
		"0",
		"1", "-1",
		"999999999999999999.999999999999999999",
		"-999999999999999999.999999999999999999",
		"0.000000000000000001",
		"-0.000000000000000001",
	}

	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			f := NewS(v)
			s := f.String()
			// Max: '-' (1) + 18 digits + '.' (1) + 18 digits = 38
			if len(s) > 38 {
				t.Errorf("String() too long (%d chars): %q", len(s), s)
			}

			sn18 := f.StringN(18)
			if sn18 != "NaN" && len(sn18) > 38 {
				t.Errorf("StringN(18) too long (%d chars): %q", len(sn18), sn18)
			}
		})
	}
}
