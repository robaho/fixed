package fixed_test

import (
	"bytes"
	"encoding/json"
	"math"
	"math/rand"
	"testing"

	. "github.com/PKartaviy/fixed"
)

func TestBasic(t *testing.T) {
	testCases := []string{
		"123.456",
		"123.456",
		"-123.456",
		"0.456",
		"-0.456",
	}

	var fs []Fixed
	for _, s := range testCases {
		f := NewS(s)
		if f.String() != s {
			t.Error("should be equal", f.String(), s)
		}
		fs = append(fs, f)
	}

	if !fs[0].Equal(fs[1]) {
		t.Error("should be equal", fs[0], fs[1])
	}

	if fs[0].Int() != 123 {
		t.Error("should be equal", fs[0].Int(), 123)
	}

	f0 := NewF(1)
	f1 := NewF(.5).Add(NewF(.5))
	f2 := NewF(.3).Add(NewF(.3)).Add(NewF(.4))

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}
	if !f0.Equal(f2) {
		t.Error("should be equal", f0, f2)
	}

	f0 = NewF(.999)
	if f0.String() != "0.999" {
		t.Error("should be equal", f0, "0.999")
	}
}

func TestNegative(t *testing.T) {
	f0 := NewS("-0.5")
	if !f0.Equal(NewF(-.5)) {
		t.Error("should be -0.5", f0)
	}
	f1 := NewS("-0.5")
	f2 := f0.Add(f1)
	if !f2.Equal(MustParse("-1")) {
		t.Error("should be -1", f2)
	}
}

func TestParse(t *testing.T) {
	_, err := Parse("123")
	if err != nil {
		t.Fail()
	}
	_, err = Parse("abc")
	if err == nil {
		t.Fail()
	}
}

func TestMustParse(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = MustParse("abc")
}

func TestNewI(t *testing.T) {
	f := NewI(123, 1)
	if f.String() != "12.3" {
		t.Error("should be equal", f, "12.3")
	}
	f = NewI(-123, 1)
	if f.String() != "-12.3" {
		t.Error("should be equal", f, "-12.3")
	}
	f = NewI(123, 0)
	if f.String() != "123" {
		t.Error("should be equal", f, "123")
	}
	f = NewI(123456789012, 9)
	if f.String() != "123.456789012" {
		t.Error("should be equal", f, "123.456789012")
	}
	f = NewI(123456789012, 9)
	if f.StringN(7) != "123.4567890" {
		t.Error("should be equal", f.StringN(7), "123.4567890")
	}

}

func TestSign(t *testing.T) {
	f0 := NewS("0")
	if f0.Sign() != 0 {
		t.Error("should be equal", f0.Sign(), 0)
	}
	f0 = NewS("NaN")
	if f0.Sign() != 0 {
		t.Error("should be equal", f0.Sign(), 0)
	}
	f0 = NewS("-100")
	if f0.Sign() != -1 {
		t.Error("should be equal", f0.Sign(), -1)
	}
	f0 = NewS("100")
	if f0.Sign() != 1 {
		t.Error("should be equal", f0.Sign(), 1)
	}

}

func TestMaxValue(t *testing.T) {
	f0 := NewS("1234567890")
	if f0.String() != "1234567890" {
		t.Error("should be equal", f0, "1234567890")
	}
	f0 = NewS("123456789012")
	if f0.String() != "123456789012" {
		t.Error("should be equal", f0, "123456789012")
	}
	f0 = NewS("-1234567890")
	if f0.String() != "-1234567890" {
		t.Error("should be equal", f0, "-1234567890")
	}
	f0 = NewS("-12345678901")
	if f0.String() != "-12345678901" {
		t.Error("should be equal", f0, "-12345678901")
	}
	f0 = NewS("9999999999")
	if f0.String() != "9999999999" {
		t.Error("should be equal", f0, "9999999999")
	}
	f0 = NewS("9.99999999")
	if f0.String() != "9.99999999" {
		t.Error("should be equal", f0, "9.99999999")
	}
	f0 = NewS("9999999999.99999999")
	if f0.String() != "9999999999.99999999" {
		t.Error("should be equal", f0, "9999999999.99999999")
	}
	f0 = NewS("9999999999.12345678901234567890")
	if f0.String() != "9999999999.123456789012345678" {
		t.Error("should be equal", f0, "9999999999.123456789012345678")
	}

}

func TestFloat(t *testing.T) {
	f0 := NewS("123.456")
	f1 := NewF(123.456)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	f1 = NewF(0.0001)

	if f1.String() != "0.0001" {
		t.Error("should be equal", f1.String(), "0.0001")
	}

	f1 = NewS(".1")
	f2 := NewS(NewF(f1.Float()).String())
	if !f1.Equal(f2) {
		t.Error("should be equal", f1, f2)
	}

}

func TestInfinite(t *testing.T) {
	f0 := NewS("0.10")
	f1 := NewF(0.10)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	f2 := NewF(0.0)
	for i := 0; i < 3; i++ {
		f2 = f2.Add(NewF(.10))
	}
	if f2.String() != "0.3" {
		t.Error("should be equal", f2.String(), "0.3")
	}

	f2 = NewF(0.0)
	for i := 0; i < 10; i++ {
		f2 = f2.Add(NewF(.10))
	}
	if f2.String() != "1" {
		t.Error("should be equal", f2.String(), "1")
	}

}

func TestAddSub(t *testing.T) {
	f0 := NewS("1")
	f1 := NewS("0.3333333")

	f2 := f0.Sub(f1)
	f2 = f2.Sub(f1)
	f2 = f2.Sub(f1)

	if f2.String() != "0.0000001" {
		t.Error("should be equal", f2.String(), "0.0000001")
	}
	f2 = f2.Sub(NewS("0.0000001"))
	if f2.String() != "0" {
		t.Error("should be equal", f2.String(), "0")
	}

	f0 = NewS("0")
	for i := 0; i < 10; i++ {
		f0 = f0.Add(NewS("0.1"))
	}
	if f0.String() != "1" {
		t.Error("should be equal", f0.String(), "1")
	}

}

func TestAbs(t *testing.T) {
	f := NewS("NaN")
	f = f.Abs()
	if !f.IsNaN() {
		t.Error("should be NaN", f)
	}
	f = NewS("1")
	f = f.Abs()
	if f.String() != "1" {
		t.Error("should be equal", f, "1")
	}
	f = NewS("-1")
	f = f.Abs()
	if f.String() != "1" {
		t.Error("should be equal", f, "1")
	}
}

func TestMulDiv(t *testing.T) {
	f0 := NewS("123.456")
	f1 := NewS("1000")

	f2 := f0.Mul(f1)
	if f2.String() != "123456" {
		t.Error("should be equal", f2.String(), "123456")
	}
	f0 = NewS("123456")
	f1 = NewS("0.0001")

	f2 = f0.Mul(f1)
	if f2.String() != "12.3456" {
		t.Error("should be equal", f2.String(), "12.3456")
	}

	f0 = NewS("123.456")
	f1 = NewS("-1000")

	f2 = f0.Mul(f1)
	if f2.String() != "-123456" {
		t.Error("should be equal", f2.String(), "-123456")
	}

	f0 = NewS("-123.456")
	f1 = NewS("-1000")

	f2 = f0.Mul(f1)
	if f2.String() != "123456" {
		t.Error("should be equal", f2.String(), "123456")
	}

	f0 = NewS("123.456")
	f1 = NewS("-1000")

	f2 = f0.Mul(f1)
	if f2.String() != "-123456" {
		t.Error("should be equal", f2.String(), "-123456")
	}

	f0 = NewS("-123.456")
	f1 = NewS("-1000")

	f2 = f0.Mul(f1)
	if f2.String() != "123456" {
		t.Error("should be equal", f2.String(), "123456")
	}

	f0 = NewS("10000.1")
	f1 = NewS("10000")

	f2 = f0.Mul(f1)
	if f2.String() != "100001000" {
		t.Error("should be equal", f2.String(), "100001000")
	}

	f2 = f2.Div(f1)
	if !f2.Equal(f0) {
		t.Error("should be equal", f0, f2)
	}

	f0 = NewS("2")
	f1 = NewS("3")

	f2 = f0.Div(f1)
	if f2.String() != "0.666666666666666667" {
		t.Error("should be equal", f2.String(), "0.666666666666666667")
	}

	f0 = NewS("1000")
	f1 = NewS("10")

	f2 = f0.Div(f1)
	if f2.String() != "100" {
		t.Error("should be equal", f2.String(), "100")
	}

	f0 = NewS("1000")
	f1 = NewS("0.1")

	f2 = f0.Div(f1)
	if f2.String() != "10000" {
		t.Error("should be equal", f2.String(), "10000")
	}

	f0 = NewS("1")
	f1 = NewS("0.1")

	f2 = f0.Mul(f1)
	if f2.String() != "0.1" {
		t.Error("should be equal", f2.String(), "0.1")
	}

	f0 = NewS("0.000001")
	f1 = NewS("0.066248")

	f2 = f0.Mul(f1)
	if f2.String() != "0.000000066248" {
		t.Error("should be equal", f2.String(), "0.000000066248")
	}

	f0 = NewS("-0.000001")
	f1 = NewS("0.066248")

	f2 = f0.Mul(f1)
	if f2.String() != "-0.000000066248" {
		t.Error("should be equal", f2.String(), "-0.000000066248")
	}

}

func TestNegatives(t *testing.T) {
	f0 := NewS("99")
	f1 := NewS("100")

	f2 := f0.Sub(f1)
	if f2.String() != "-1" {
		t.Error("should be equal", f2.String(), "-1")
	}
	f0 = NewS("-1")
	f1 = NewS("-1")

	f2 = f0.Sub(f1)
	if f2.String() != "0" {
		t.Error("should be equal", f2.String(), "0")
	}
	f0 = NewS(".001")
	f1 = NewS(".002")

	f2 = f0.Sub(f1)
	if f2.String() != "-0.001" {
		t.Error("should be equal", f2.String(), "-0.001")
	}
}

func TestOverflow(t *testing.T) {
	f0 := NewF(1.1234567)
	if f0.String() != "1.1234567" {
		t.Error("should be equal", f0.String(), "1.1234567")
	}
	f0 = NewF(1.123456789123)
	if f0.String() != "1.123456789123" {
		t.Error("should be equal", f0.String(), "1.123456789123")
	}
	f0 = NewF(1.0 / 3.0)
	if f0.String() != "0.3333333333333333" {
		t.Error("should be equal", f0.String(), "0.3333333333333333")
	}
	f0 = NewF(2.0 / 3.0)
	if f0.String() != "0.6666666666666666" {
		t.Error("should be equal", f0.String(), "0.6666666666666666")
	}

}

func TestNaN(t *testing.T) {
	f0 := NewF(math.NaN())
	if !f0.IsNaN() {
		t.Error("f0 should be NaN")
	}
	if f0.String() != "NaN" {
		t.Error("should be equal", f0.String(), "NaN")
	}
	f0 = NewS("NaN")
	if !f0.IsNaN() {
		t.Error("f0 should be NaN")
	}

	f0 = NewS("0.0004096")
	if f0.String() != "0.0004096" {
		t.Error("should be equal", f0.String(), "0.0004096")
	}

}

func TestIntFrac(t *testing.T) {
	f0 := NewF(1234.5678)
	if f0.Int() != 1234 {
		t.Error("should be equal", f0.Int(), 1234)
	}
	if f0.Frac() != .5678 {
		t.Error("should be equal", f0.Frac(), .5678)
	}
}

func TestString(t *testing.T) {
	f0 := NewF(1234.5678)
	if f0.String() != "1234.5678" {
		t.Error("should be equal", f0.String(), "1234.5678")
	}
	f0 = NewF(1234.0)
	if f0.String() != "1234" {
		t.Error("should be equal", f0.String(), "1234")
	}
}

func TestStringN(t *testing.T) {
	f0 := NewS("1.1")
	s := f0.StringN(2)

	if s != "1.10" {
		t.Error("should be equal", s, "1.10")
	}
	f0 = NewS("1")
	s = f0.StringN(2)

	if s != "1.00" {
		t.Error("should be equal", s, "1.00")
	}

	f0 = NewS("1.123")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}
	f0 = NewS("1.123")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}

	f0 = NewS("1.127")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}

	f0 = NewS("1.123")
	s = f0.StringN(0)

	if s != "1" {
		t.Error("should be equal", s, "1")
	}
}

func TestRound(t *testing.T) {
	f0 := NewS("1.12345")
	f1 := f0.Round(2)

	if f1.String() != "1.12" {
		t.Error("should be equal", f1, "1.12")
	}

	f1 = f0.Round(5)

	if f1.String() != "1.12345" {
		t.Error("should be equal", f1, "1.12345")
	}
	f1 = f0.Round(4)

	if f1.String() != "1.1235" {
		t.Error("should be equal", f1, "1.1235")
	}

	f0 = NewS("-1.12345")
	f1 = f0.Round(3)

	if f1.String() != "-1.123" {
		t.Error("should be equal", f1, "-1.123")
	}
	f0 = NewS("-1.1235")
	f1 = f0.Round(3)

	if f1.String() != "-1.124" {
		t.Error("should be equal", f1, "-1.124")
	}
	f1 = f0.Round(4)

	if f1.String() != "-1.1235" {
		t.Error("should be equal", f1, "-1.1235")
	}

	f0 = NewS("-0.0001")
	f1 = f0.Round(1)

	if f1.String() != "0" {
		t.Error("should be equal", f1, "0")
	}

	f0 = NewS("2234.565")
	f1 = f0.Round(2)

	if f1.String() != "2234.57" {
		t.Error("should be equal", f1, "2234.57")
	}

	f0 = NewS("123.456")
	f1 = f0.Round(-1)

	if f1.String() != "120" {
		t.Error("should be equal", f1, "120")
	}

	f1 = f0.Round(-2)

	if f1.String() != "100" {
		t.Error("should be equal", f1, "100")
	}

	f0 = NewS("125.456")
	f1 = f0.Round(-1)

	if f1.String() != "130" {
		t.Error("should be equal", f1, "130")
	}
}

func TestFloorCeil(t *testing.T) {
	f := NewS("1.5")
	if f.Floor(0).String() != "1" {
		t.Error("Floor(1.5, 0) should be 1, got", f.Floor(0))
	}
	if f.Ceil(0).String() != "2" {
		t.Error("Ceil(1.5, 0) should be 2, got", f.Ceil(0))
	}

	f = NewS("1.234")
	if f.Floor(2).String() != "1.23" {
		t.Error("Floor(1.234, 2) should be 1.23, got", f.Floor(2))
	}
	if f.Ceil(2).String() != "1.24" {
		t.Error("Ceil(1.234, 2) should be 1.24, got", f.Ceil(2))
	}

	f = NewS("-1.5")
	if f.Floor(0).String() != "-2" {
		t.Error("Floor(-1.5, 0) should be -2, got", f.Floor(0))
	}
	if f.Ceil(0).String() != "-1" {
		t.Error("Ceil(-1.5, 0) should be -1, got", f.Ceil(0))
	}

	f = NewS("-1.234")
	if f.Floor(2).String() != "-1.24" {
		t.Error("Floor(-1.234, 2) should be -1.24, got", f.Floor(2))
	}
	if f.Ceil(2).String() != "-1.23" {
		t.Error("Ceil(-1.234, 2) should be -1.23, got", f.Ceil(2))
	}

	f = NewS("123.456")
	if f.Floor(-1).String() != "120" {
		t.Error("Floor(123.456, -1) should be 120, got", f.Floor(-1))
	}
	if f.Ceil(-1).String() != "130" {
		t.Error("Ceil(123.456, -1) should be 130, got", f.Ceil(-1))
	}

	f = NewS("1.2")
	if f.Floor(1).String() != "1.2" {
		t.Error("Floor(1.2, 1) should be 1.2, got", f.Floor(1))
	}
	if f.Ceil(1).String() != "1.2" {
		t.Error("Ceil(1.2, 1) should be 1.2, got", f.Ceil(1))
	}

	f = NewS("0")
	if f.Floor(0).String() != "0" {
		t.Error("Floor(0, 0) should be 0, got", f.Floor(0))
	}
	if f.Ceil(0).String() != "0" {
		t.Error("Ceil(0, 0) should be 0, got", f.Ceil(0))
	}

	f = NewS("NaN")
	if !f.Floor(0).IsNaN() {
		t.Error("Floor(NaN, 0) should be NaN, got", f.Floor(0))
	}
	if !f.Ceil(0).IsNaN() {
		t.Error("Ceil(NaN, 0) should be NaN, got", f.Ceil(0))
	}
}

func TestEncodeDecode(t *testing.T) {
	b := &bytes.Buffer{}

	f := NewS("12345.12345")

	f.WriteTo(b)

	f0, err := ReadFrom(b)
	if err != nil {
		t.Error(err)
	}

	if !f.Equal(f0) {
		t.Error("don't match", f, f0)
	}

	data, err := f.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	f1 := NewF(0)
	f1.UnmarshalBinary(data)

	if !f.Equal(f1) {
		t.Error("don't match", f, f0)
	}
}

type JStruct struct {
	F Fixed `json:"f"`
}

func TestJSON(t *testing.T) {
	j := JStruct{}

	f := NewS("1234567.1234567")
	j.F = f

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	err := enc.Encode(&j)
	if err != nil {
		t.Error(err)
	}

	j.F = ZERO

	dec := json.NewDecoder(&buf)

	err = dec.Decode(&j)
	if err != nil {
		t.Error(err)
	}

	if !j.F.Equal(f) {
		t.Error("don't match", j.F, f)
	}
}

func TestJSON_NaN(t *testing.T) {
	j := JStruct{}

	f := NewS("NaN")
	j.F = f

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	err := enc.Encode(&j)
	if err != nil {
		t.Error(err)
	}

	j.F = ZERO

	dec := json.NewDecoder(&buf)

	err = dec.Decode(&j)
	if err != nil {
		t.Error(err)
	}

	if !j.F.IsNaN() {
		t.Error("did not decode NaN", j.F, f)
	}
}

func TestMulVsMulSlow(t *testing.T) {
	type pair struct {
		name string
		a, b Fixed
	}

	one := NewS("1")
	negOne := NewS("-1")
	zero := NewS("0")

	cases := []pair{
		{"pos*pos", NewS("123.456"), NewS("789.012")},
		{"pos*neg", NewS("123.456"), NewS("-789.012")},
		{"neg*neg", NewS("-123.456"), NewS("-789.012")},
		{"zero*pos", zero, NewS("123.456")},
		{"pos*zero", NewS("123.456"), zero},
		{"zero*neg", zero, NewS("-123.456")},
		{"neg*zero", NewS("-123.456"), zero},
		{"zero*zero", zero, zero},
		{"one*val", one, NewS("999.999")},
		{"negone*val", negOne, NewS("999.999")},
		{"frac*frac", NewS("0.000001"), NewS("0.066248")},
		{"neg_frac*frac", NewS("-0.000001"), NewS("0.066248")},
		{"large_int", NewS("10000.1"), NewS("10000")},
		{"near_max", NewS("999999999"), NewS("999999999")},
		{"small_fracs", NewS("0.000000000000000001"), NewS("1")},
		{"mixed1", NewS("123456789.123456789"), NewS("0.000000001")},
		{"mixed2", NewS("1.999999999999999999"), NewS("1.999999999999999999")},
		// Large hi values (a[0] > 0 in base-10^9 decomposition)
		{"large_hi*frac", NewS("1000000001.5"), NewS("0.5")},
		{"large_hi*one", NewS("999999999999999999"), NewS("1")},
		{"large_hi*small", NewS("999999999999999999"), NewS("0.000000000000000001")},
		{"large_hi*neg_frac", NewS("999999999999999999"), NewS("-0.000000000000000001")},
		{"large_hi*large_frac", NewS("123456789012345678"), NewS("0.000000001")},
		{"neg_large*frac", NewS("-123456789012345678"), NewS("0.5")},
		// Overflow cases
		{"overflow", NewS("999999999999999999"), NewS("999999999999999999")},
		{"near_overflow", NewS("999999999999999999.999999999999999999"), NewS("1.000000000000000001")},
		{"both_nan", NaN, NaN},
		{"nan_left", NaN, one},
		{"nan_right", one, NaN},
	}

	for _, tc := range cases {
		fast := tc.a.Mul(tc.b)
		slow := tc.a.MulSlow(tc.b)
		if fast.IsNaN() && slow.IsNaN() {
			continue
		}
		if !fast.Equal(slow) {
			t.Errorf("%s: Mul(%s, %s) = %s, MulSlow = %s",
				tc.name, tc.a, tc.b, fast, slow)
		}
	}

	// Random/fuzz loop with two ranges:
	// 1) Small hi (products usually in range) — tests normal paths
	// 2) Large hi with small multiplier — tests a[0]>0 digit decomposition
	rng := rand.New(rand.NewSource(42))

	compare := func(tag string, i int, a, b Fixed) {
		fast := a.Mul(b)
		slow := a.MulSlow(b)
		if fast.IsNaN() && slow.IsNaN() {
			return
		}
		if !fast.Equal(slow) {
			t.Errorf("%s iter %d: Mul(%s, %s) = %s, MulSlow = %s",
				tag, i, a, b, fast, slow)
		}
	}

	for i := 0; i < 10000; i++ {
		// Small hi: [-10^9, 10^9]
		aHi := rng.Int63n(2000000000) - 1000000000
		aLo := rng.Int63n(1000000000000000000)
		bHi := rng.Int63n(2000000000) - 1000000000
		bLo := rng.Int63n(1000000000000000000)
		if aHi < 0 {
			aLo = -aLo
		}
		if bHi < 0 {
			bLo = -bLo
		}
		a := NewI(aHi, 0).Add(NewI(aLo, 18))
		b := NewI(bHi, 0).Add(NewI(bLo, 18))
		compare("small", i, a, b)
	}

	for i := 0; i < 10000; i++ {
		// Large hi * small value — exercises a[0] > 0 without always overflowing
		aHi := rng.Int63n(999999999999999999) + 1 // [1, maxHi]
		aLo := rng.Int63n(1000000000000000000)
		// Small multiplier: pure fraction in [0, 1)
		bLo := rng.Int63n(1000000000000000000)
		sign := int64(1)
		if rng.Intn(2) == 0 {
			sign = -1
		}
		a := NewI(sign*aHi, 0).Add(NewI(sign*aLo, 18))
		b := NewI(bLo, 18)
		compare("large", i, a, b)
	}
}

func TestDivVsDivSlow(t *testing.T) {
	type pair struct {
		name string
		a, b Fixed
	}

	one := NewS("1")
	negOne := NewS("-1")
	zero := NewS("0")

	cases := []pair{
		{"1/3", one, NewS("3")},
		{"2/3", NewS("2"), NewS("3")},
		{"1/7", one, NewS("7")},
		{"1/9", one, NewS("9")},
		{"10/3", NewS("10"), NewS("3")},
		{"-1/3", negOne, NewS("3")},
		{"1/-3", one, NewS("-3")},
		{"-1/-3", negOne, NewS("-3")},
		{"5/5", NewS("5"), NewS("5")},
		{"large/small", NewS("999999999999999999"), NewS("0.000000000000000001")},
		{"small/large", NewS("0.000000000000000001"), NewS("999999999999999999")},
		{"max/1", NewS("999999999999999999.999999999999999999"), one},
		{"max/-1", NewS("999999999999999999.999999999999999999"), negOne},
		{"frac/frac", NewS("0.123456789"), NewS("0.987654321")},
		{"-frac/frac", NewS("-0.123456789"), NewS("0.987654321")},
		{"large_int/large_int", NewS("123456789012345678"), NewS("987654321")},
		{"1/1", one, one},
		{"0/1", zero, one},
		{"0/-1", zero, negOne},
		{"div_by_zero", one, zero},
		{"mixed", NewS("123456789.123456789"), NewS("0.000000001")},
		{"near_max_result", NewS("999999999999999999"), NewS("1.000000000000000001")},
		{"both_nan", NaN, NaN},
		{"nan_left", NaN, one},
		{"nan_right", one, NaN},
	}

	for _, tc := range cases {
		fast := tc.a.Div(tc.b)
		slow := tc.a.DivSlow(tc.b)
		if fast.IsNaN() && slow.IsNaN() {
			continue
		}
		if !fast.Equal(slow) {
			t.Errorf("%s: Div(%s, %s) = %s, DivSlow = %s",
				tc.name, tc.a, tc.b, fast, slow)
		}
	}

	// Random/fuzz loop
	rng := rand.New(rand.NewSource(42))

	compare := func(tag string, i int, a, b Fixed) {
		fast := a.Div(b)
		slow := a.DivSlow(b)
		if fast.IsNaN() && slow.IsNaN() {
			return
		}
		if !fast.Equal(slow) {
			t.Errorf("%s iter %d: Div(%s, %s) = %s, DivSlow = %s",
				tag, i, a, b, fast, slow)
		}
	}

	for i := 0; i < 10000; i++ {
		// Random values with moderate hi
		aHi := rng.Int63n(2000000000) - 1000000000
		aLo := rng.Int63n(1000000000000000000)
		bHi := rng.Int63n(2000000000) - 1000000000
		bLo := rng.Int63n(1000000000000000000)
		if aHi < 0 {
			aLo = -aLo
		}
		if bHi < 0 {
			bLo = -bLo
		}
		a := NewI(aHi, 0).Add(NewI(aLo, 18))
		b := NewI(bHi, 0).Add(NewI(bLo, 18))
		// Skip division by zero
		if b.Sign() == 0 {
			continue
		}
		compare("small", i, a, b)
	}

	for i := 0; i < 10000; i++ {
		// Large dividend / small divisor — exercises quotient overflow paths
		aHi := rng.Int63n(999999999999999999) + 1
		aLo := rng.Int63n(1000000000000000000)
		// Divisor: small positive value to avoid trivial overflow
		bHi := rng.Int63n(1000) + 1
		bLo := rng.Int63n(1000000000000000000)
		sign := int64(1)
		if rng.Intn(2) == 0 {
			sign = -1
		}
		a := NewI(sign*aHi, 0).Add(NewI(sign*aLo, 18))
		b := NewI(bHi, 0).Add(NewI(bLo, 18))
		compare("large", i, a, b)
	}
}

func TestNewFVsNewFSlow(t *testing.T) {
	specific := []float64{
		0, -0.0, 1, -1, 0.5, -0.5,
		123.456, -123.456,
		0.0001, -0.0001,
		1.0 / 3.0, 2.0 / 3.0,
		math.Pi, math.E,
		1e15, -1e15,
		999999999999999999.0,
		-999999999999999999.0,
		0.000000000000000001,
		42.0,
		0.1, 0.2, 0.3,
		1234567890.123456789,
	}

	match := func(label string, v float64) {
		fast := NewF(v)
		slow := NewFSlow(v)
		// NaN.Equal(NaN) is false, so compare IsNaN separately
		if fast.IsNaN() && slow.IsNaN() {
			return
		}
		if !fast.Equal(slow) {
			t.Errorf("%s: NewF(%v) = %s, NewFSlow(%v) = %s", label, v, fast.String(), v, slow.String())
		}
	}

	for _, v := range specific {
		match("specific", v)
	}

	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 10000; i++ {
		f := (rng.Float64()*2 - 1) * 1e18
		if f >= MAX || f <= -MAX {
			continue
		}
		match("random", f)
	}
}
