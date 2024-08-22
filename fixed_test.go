package fixed_test

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"

	. "github.com/robaho/fixed"
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
	if f.String() != "123.456789" {
		t.Error("should be equal", f, "123.456789")
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
	f0 := NewS("12345678901")
	if f0.String() != "12345678901" {
		t.Error("should be equal", f0, "12345678901")
	}
	f0 = NewS("123456789012")
	if f0.String() != "NaN" {
		t.Error("should be equal", f0, "NaN")
	}
	f0 = NewS("-12345678901")
	if f0.String() != "-12345678901" {
		t.Error("should be equal", f0, "-12345678901")
	}
	f0 = NewS("-123456789012")
	if f0.String() != "NaN" {
		t.Error("should be equal", f0, "NaN")
	}
	f0 = NewS("99999999999")
	if f0.String() != "99999999999" {
		t.Error("should be equal", f0, "99999999999")
	}
	f0 = NewS("9.9999999")
	if f0.String() != "9.9999999" {
		t.Error("should be equal", f0, "9.9999999")
	}
	f0 = NewS("99999999999.9999999")
	if f0.String() != "99999999999.9999999" {
		t.Error("should be equal", f0, "99999999999.9999999")
	}
	f0 = NewS("99999999999.12345678901234567890")
	if f0.String() != "99999999999.1234567" {
		t.Error("should be equal", f0, "99999999999.1234567")
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
	if f2.String() != "0.6666667" {
		t.Error("should be equal", f2.String(), "0.6666667")
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
	if f2.String() != "0.0000001" {
		t.Error("should be equal", f2.String(), "0.0000001")
	}

	f0 = NewS("-0.000001")
	f1 = NewS("0.066248")

	f2 = f0.Mul(f1)
	if f2.String() != "-0.0000001" {
		t.Error("should be equal", f2.String(), "-0.0000001")
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
	if f0.String() != "1.1234568" {
		t.Error("should be equal", f0.String(), "1.1234568")
	}
	f0 = NewF(1.0 / 3.0)
	if f0.String() != "0.3333333" {
		t.Error("should be equal", f0.String(), "0.3333333")
	}
	f0 = NewF(2.0 / 3.0)
	if f0.String() != "0.6666667" {
		t.Error("should be equal", f0.String(), "0.6666667")
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
