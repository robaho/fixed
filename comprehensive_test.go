package fixed_test

import (
	"math"
	"testing"

	. "github.com/PKartaviy/fixed"
)

// TestSignConsistency ensures hi and lo always have consistent signs
func TestSignConsistency(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"positive", "123.456"},
		{"negative", "-123.456"},
		{"negative fraction", "-0.456"},
		{"positive fraction", "0.456"},
		{"negative small", "-0.000000000000000001"},
		{"positive small", "0.000000000000000001"},
		{"negative large", "-999999999999999999"},
		{"positive large", "999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewS(tt.input)
			if f.IsNaN() {
				t.Fatalf("failed to parse %s", tt.input)
			}

			// Verify round-trip
			if f.String() != tt.input {
				t.Errorf("round-trip failed: got %s, want %s", f.String(), tt.input)
			}

			// Verify sign consistency by checking operations
			doubled := f.Add(f)
			mulByTwo := f.Mul(NewI(2, 0))

			// Both should be valid or both should be NaN (overflow)
			if doubled.IsNaN() != mulByTwo.IsNaN() {
				t.Errorf("sign inconsistency detected: add vs mul differ (NaN mismatch)")
			} else if !doubled.IsNaN() && !doubled.Equal(mulByTwo) {
				t.Errorf("sign inconsistency detected: add vs mul differ (values: %s vs %s)", doubled.String(), mulByTwo.String())
			}
		})
	}
}

// TestNegativeZeroHandling tests the special case of -0.xxx
func TestNegativeZeroHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-0.5", "-0.5"},
		{"-0.000000000000000001", "-0.000000000000000001"},
		{"-0.999999999999999999", "-0.999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := NewS(tt.input)
			if f.String() != tt.expected {
				t.Errorf("got %s, want %s", f.String(), tt.expected)
			}

			// Verify sign
			if f.Sign() != -1 {
				t.Errorf("sign should be -1, got %d", f.Sign())
			}

			// Verify operations preserve sign
			doubled := f.Add(f)
			if doubled.Sign() != -1 {
				t.Errorf("doubling negative should stay negative")
			}
		})
	}
}

// TestArithmeticOverflow tests operations that exceed the valid range
func TestArithmeticOverflow(t *testing.T) {
	maxVal := NewS("999999999999999999.999999999999999999")
	one := NewI(1, 0)

	// Addition overflow
	result := maxVal.Add(one)
	if !result.IsNaN() {
		t.Error("addition overflow should produce NaN")
	}

	// Subtraction underflow
	minVal := NewS("-999999999999999999.999999999999999999")
	result = minVal.Sub(one)
	if !result.IsNaN() {
		t.Error("subtraction underflow should produce NaN")
	}

	// Multiplication overflow
	large := NewS("1000000000000.0")
	result = large.Mul(large)
	if !result.IsNaN() {
		t.Error("multiplication overflow should produce NaN")
	}
}

// TestMultiplicationAccuracy validates multiplication precision
func TestMultiplicationAccuracy(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"123.456", "789.012", "97408.265472"},
		{"0.000000000000000001", "999999999999999999", "0.999999999999999999"},
		{"999999999999999999", "0.000000000000000001", "0.999999999999999999"},
		{"-123.456", "789.012", "-97408.265472"},
		{"-123.456", "-789.012", "97408.265472"},
		{"0.333333333333333333", "3", "0.999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.a+"*"+tt.b, func(t *testing.T) {
			a := NewS(tt.a)
			b := NewS(tt.b)
			expected := NewS(tt.expected)

			result := a.Mul(b)
			if !result.Equal(expected) {
				t.Errorf("got %s, want %s", result.String(), expected.String())
			}
		})
	}
}

// TestDivisionAccuracy validates division precision
func TestDivisionAccuracy(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"1", "3", "0.333333333333333333"},
		{"2", "3", "0.666666666666666667"},
		{"1", "7", "0.142857142857142857"},
		{"1", "9", "0.111111111111111111"},
		{"10", "3", "3.333333333333333333"},
		{"-1", "3", "-0.333333333333333333"},
		{"1", "-3", "-0.333333333333333333"},
		{"-1", "-3", "0.333333333333333333"},
	}

	for _, tt := range tests {
		t.Run(tt.a+"/"+tt.b, func(t *testing.T) {
			a := NewS(tt.a)
			b := NewS(tt.b)
			expected := NewS(tt.expected)

			result := a.Div(b)
			if !result.Equal(expected) {
				t.Errorf("got %s, want %s", result.String(), expected.String())
			}
		})
	}
}

// TestDivisionByZero validates NaN handling
func TestDivisionByZero(t *testing.T) {
	one := NewI(1, 0)
	zero := NewI(0, 0)

	result := one.Div(zero)
	if !result.IsNaN() {
		t.Error("division by zero should produce NaN")
	}
}

// TestNaNPropagation ensures NaN propagates through operations
func TestNaNPropagation(t *testing.T) {
	nan := NaN
	one := NewI(1, 0)

	tests := []struct {
		name   string
		result Fixed
	}{
		{"add", nan.Add(one)},
		{"sub", nan.Sub(one)},
		{"mul", nan.Mul(one)},
		{"div", nan.Div(one)},
		{"add reverse", one.Add(nan)},
		{"sub reverse", one.Sub(nan)},
		{"mul reverse", one.Mul(nan)},
		{"div reverse", one.Div(nan)},
		{"abs", nan.Abs()},
		{"round", nan.Round(2)},
		{"floor", nan.Floor(2)},
		{"ceil", nan.Ceil(2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.result.IsNaN() {
				t.Error("NaN should propagate through operation")
			}
		})
	}
}

// TestRoundingEdgeCases tests rounding at various decimal places
func TestRoundingEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		places   int
		expected string
	}{
		{"1.5", 0, "2"},
		{"1.4", 0, "1"},
		{"-1.5", 0, "-2"},
		{"-1.4", 0, "-1"},
		{"1.555555555555555555", 1, "1.6"},
		{"1.444444444444444444", 1, "1.4"},
		{"123.456789012345678", 10, "123.4567890123"},
		{"123.456789012345678", 18, "123.456789012345678"},
		{"999.999999999999999", 0, "1000"},
		{"-999.999999999999999", 0, "-1000"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := NewS(tt.input)
			result := f.Round(tt.places)
			expected := NewS(tt.expected)

			if !result.Equal(expected) {
				t.Errorf("Round(%d): got %s, want %s", tt.places, result.String(), expected.String())
			}
		})
	}
}

// TestFloorCeilEdgeCases tests floor and ceil operations
func TestFloorCeilEdgeCases(t *testing.T) {
	tests := []struct {
		input         string
		places        int
		expectedFloor string
		expectedCeil  string
	}{
		{"1.9", 0, "1", "2"},
		{"1.1", 0, "1", "2"},
		{"-1.9", 0, "-2", "-1"},
		{"-1.1", 0, "-2", "-1"},
		{"0.1", 0, "0", "1"},
		{"-0.1", 0, "-1", "0"},
		{"123.456789012345678", 5, "123.45678", "123.45679"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := NewS(tt.input)

			floor := f.Floor(tt.places)
			expectedFloor := NewS(tt.expectedFloor)
			if !floor.Equal(expectedFloor) {
				t.Errorf("Floor(%d): got %s, want %s", tt.places, floor.String(), expectedFloor.String())
			}

			ceil := f.Ceil(tt.places)
			expectedCeil := NewS(tt.expectedCeil)
			if !ceil.Equal(expectedCeil) {
				t.Errorf("Ceil(%d): got %s, want %s", tt.places, ceil.String(), expectedCeil.String())
			}
		})
	}
}

// TestComparisonEdgeCases tests comparison operations
func TestComparisonEdgeCases(t *testing.T) {
	tests := []struct {
		a, b string
		cmp  int
	}{
		{"0", "0", 0},
		{"1", "0", 1},
		{"0", "1", -1},
		{"-1", "0", -1},
		{"0", "-1", 1},
		{"0.000000000000000001", "0", 1},
		{"0", "0.000000000000000001", -1},
		{"-0.000000000000000001", "0", -1},
		{"0", "-0.000000000000000001", 1},
		{"999999999999999999.999999999999999999", "999999999999999999.999999999999999998", 1},
		{"-999999999999999999.999999999999999999", "-999999999999999999.999999999999999998", -1},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			a := NewS(tt.a)
			b := NewS(tt.b)

			cmp := a.Cmp(b)
			if cmp != tt.cmp {
				t.Errorf("Cmp: got %d, want %d", cmp, tt.cmp)
			}

			// Test Equal consistency
			if tt.cmp == 0 && !a.Equal(b) {
				t.Error("Equal should return true when Cmp returns 0")
			}
			if tt.cmp != 0 && a.Equal(b) {
				t.Error("Equal should return false when Cmp returns non-zero")
			}
		})
	}
}

// TestBinarySerializationRoundTrip tests serialization edge cases
func TestBinarySerializationRoundTrip(t *testing.T) {
	tests := []string{
		"0",
		"1",
		"-1",
		"0.000000000000000001",
		"-0.000000000000000001",
		"999999999999999999.999999999999999999",
		"-999999999999999999.999999999999999999",
		"123.456789012345678",
		"-123.456789012345678",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			original := NewS(tt)

			data, err := original.MarshalBinary()
			if err != nil {
				t.Fatalf("MarshalBinary failed: %v", err)
			}

			var restored Fixed
			err = restored.UnmarshalBinary(data)
			if err != nil {
				t.Fatalf("UnmarshalBinary failed: %v", err)
			}

			if !original.Equal(restored) {
				t.Errorf("round-trip failed: got %s, want %s", restored.String(), original.String())
			}
		})
	}

	// Test NaN serialization
	t.Run("NaN", func(t *testing.T) {
		data, err := NaN.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary failed: %v", err)
		}

		var restored Fixed
		err = restored.UnmarshalBinary(data)
		if err != nil {
			t.Fatalf("UnmarshalBinary failed: %v", err)
		}

		if !restored.IsNaN() {
			t.Error("NaN should serialize and deserialize as NaN")
		}
	})
}

// TestStringParsingEdgeCases tests string parsing edge cases
func TestStringParsingEdgeCases(t *testing.T) {
	tests := []struct {
		input       string
		shouldParse bool
		expected    string
	}{
		{"0", true, "0"},
		{".5", true, "0.5"},
		{"-.5", true, "-0.5"},
		{"5.", true, "5"},
		{"-5.", true, "-5"},
		{"00123.45", true, "123.45"},
		{"-00123.45", true, "-123.45"},
		{"123.450000000000000000", true, "123.45"},
		{"+123", true, "123"},
		{"", false, ""},
		{"abc", false, ""},
		{"12.34.56", false, ""},
		{"--123", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f, err := NewSErr(tt.input)

			if tt.shouldParse {
				if err != nil {
					t.Errorf("should parse successfully, got error: %v", err)
				}
				if f.String() != tt.expected {
					t.Errorf("got %s, want %s", f.String(), tt.expected)
				}
			} else {
				if err == nil {
					t.Error("should fail to parse, but succeeded")
				}
			}
		})
	}
}

// TestOperationCommutativity tests commutative properties
func TestOperationCommutativity(t *testing.T) {
	tests := []struct {
		a, b string
	}{
		{"123.456", "789.012"},
		{"0.000000000000000001", "999999999999999999"},
		{"-123.456", "789.012"},
		{"-123.456", "-789.012"},
	}

	for _, tt := range tests {
		t.Run(tt.a+","+tt.b, func(t *testing.T) {
			a := NewS(tt.a)
			b := NewS(tt.b)

			// Addition commutativity
			if !a.Add(b).Equal(b.Add(a)) {
				t.Error("addition is not commutative")
			}

			// Multiplication commutativity
			if !a.Mul(b).Equal(b.Mul(a)) {
				t.Error("multiplication is not commutative")
			}
		})
	}
}

// TestOperationAssociativity tests associative properties
func TestOperationAssociativity(t *testing.T) {
	a := NewS("123.456")
	b := NewS("789.012")
	c := NewS("345.678")

	// Addition associativity: (a + b) + c = a + (b + c)
	left := a.Add(b).Add(c)
	right := a.Add(b.Add(c))
	if !left.Equal(right) {
		t.Errorf("addition is not associative: %s != %s", left.String(), right.String())
	}

	// Multiplication associativity: (a * b) * c = a * (b * c)
	left = a.Mul(b).Mul(c)
	right = a.Mul(b.Mul(c))
	if !left.Equal(right) {
		t.Errorf("multiplication is not associative: %s != %s", left.String(), right.String())
	}
}

// TestOperationIdentities tests identity elements
func TestOperationIdentities(t *testing.T) {
	values := []string{
		"123.456",
		"-123.456",
		"0.000000000000000001",
		"999999999999999999",
	}

	zero := NewI(0, 0)
	one := NewI(1, 0)

	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			f := NewS(v)

			// Addition identity: a + 0 = a
			if !f.Add(zero).Equal(f) {
				t.Error("addition identity failed")
			}

			// Multiplication identity: a * 1 = a
			if !f.Mul(one).Equal(f) {
				t.Error("multiplication identity failed")
			}

			// Division identity: a / 1 = a
			if !f.Div(one).Equal(f) {
				t.Error("division identity failed")
			}
		})
	}
}

// TestFloatConversionPrecision tests float64 conversion limits
func TestFloatConversionPrecision(t *testing.T) {
	tests := []float64{
		0,
		1,
		-1,
		0.5,
		-0.5,
		123.456,
		-123.456,
		1.0 / 3.0,
		2.0 / 3.0,
		math.Pi,
		math.E,
		1e15,  // Near float64 precision limit
		-1e15,
	}

	for _, v := range tests {
		t.Run(formatFloat(v), func(t *testing.T) {
			f := NewF(v)

			if f.IsNaN() {
				t.Error("valid float converted to NaN")
			}

			// Convert back and check relative error
			back := f.Float()
			relError := math.Abs((back - v) / v)

			// Allow small relative error due to float64 precision
			if v != 0 && relError > 1e-15 {
				t.Errorf("float round-trip error too large: %.20f vs %.20f (rel error: %e)", back, v, relError)
			}
		})
	}
}

// TestMaxValueHandling tests operations near maximum values
func TestMaxValueHandling(t *testing.T) {
	// Test values just below max
	almostMax := NewS("999999999999999999")
	one := NewI(1, 0)
	tiny := NewS("0.000000000000000001")

	// Should succeed
	result := almostMax.Add(tiny)
	if result.IsNaN() {
		t.Error("adding tiny to almost-max should not overflow")
	}

	// Should overflow
	result = almostMax.Add(one)
	if !result.IsNaN() {
		t.Error("adding 1 to almost-max should overflow")
	}
}

// TestIntegerOperations tests operations on integer values
func TestIntegerOperations(t *testing.T) {
	a := NewI(100, 0)
	b := NewI(50, 0)

	if a.Add(b).String() != "150" {
		t.Errorf("integer addition failed: got %s", a.Add(b).String())
	}

	if a.Sub(b).String() != "50" {
		t.Errorf("integer subtraction failed: got %s", a.Sub(b).String())
	}

	if a.Mul(b).String() != "5000" {
		t.Errorf("integer multiplication failed: got %s", a.Mul(b).String())
	}

	if a.Div(b).String() != "2" {
		t.Errorf("integer division failed: got %s", a.Div(b).String())
	}
}

func formatFloat(f float64) string {
	return NewF(f).String()
}
