package fixed

// release under the terms of file license.txt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Fixed is a fixed precision 38.24 number (supports 11.7 digits). It supports NaN.
type Fixed struct {
	fp int64
}

const nPlaces = 7
const pow7 = int64(10 * 10 * 10 * 10 * 10 * 10 * 10)
const zeros = "0000000"
const nan = int64(1 << 62)

var NaN = Fixed{fp: nan}
var ZERO = Fixed{fp: 0}

const MAX = float64(99999999999.9999999)

var errTooLarge = errors.New("significand too large")
var errFormat = errors.New("invalid encoding")

// NewS creates a new Fixed from a string, returning NaN if the string could not be parsed
func NewS(s string) Fixed {
	f, _ := NewSErr(s)
	return f
}

// NewSErr creates a new Fixed from a string, returning NaN, and error if the string could not be parsed
func NewSErr(s string) (Fixed, error) {
	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NaN, err
		}
		return NewF(f), nil
	}
	if "NaN" == s {
		return NaN, nil
	}
	period := strings.Index(s, ".")
	var i int64
	var f int64
	var sign int64 = 1
	if period == -1 {
		i, _ = strconv.ParseInt(s, 10, 64)
		if i < 0 {
			sign = -1
			i = i * -1
		}
	} else {
		i, _ = strconv.ParseInt(s[:period], 10, 64)
		if i < 0 {
			sign = -1
			i = i * -1
		}
		fs := s[period+1:]
		fs = fs + zeros[:max(0, nPlaces-len(fs))]
		f, _ = strconv.ParseInt(fs[0:nPlaces], 10, 64)
	}
	if i > 99999999999 {
		return NaN, errTooLarge
	}
	return Fixed{fp: sign * (i*pow7 + f)}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewF creates a Fixed from an float64, rounding at the 8th decimal place
func NewF(f float64) Fixed {
	if math.IsNaN(f) {
		return Fixed{fp: nan}
	}
	if f >= MAX || f <= -MAX {
		return NaN
	}
	round := .5
	if f < 0 {
		round = -0.5
	}

	return Fixed{fp: int64(f*float64(pow7) + round)}
}

// NewI creates a Fixed for an integer, moving the decimal point n places to the left
// For example, NewI(123,1) becomes 12.3. If n > 7, the value is truncated
func NewI(i int64, n uint) Fixed {
	if n > nPlaces {
		i = i / int64(math.Pow10(int(n-nPlaces)))
		n = nPlaces
	}

	i = i * int64(math.Pow10(int(nPlaces-n)))

	return Fixed{fp: i}
}

func (f Fixed) IsNaN() bool {
	return f.fp == nan
}

func (f Fixed) IsZero() bool {
	return f.Equal(ZERO)
}

// Sign returns:
//
//	-1 if f <  0
//	 0 if f == 0 or NaN
//	+1 if f >  0
//
func (f Fixed) Sign() int {
	if f.IsNaN() {
		return 0
	}
	return f.Cmp(ZERO)
}

// Float converts the Fixed to a float64
func (f Fixed) Float() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp) / float64(pow7)
}

// Add adds f0 to f producing a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Add(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return Fixed{fp: f.fp + f0.fp}
}

// Sub subtracts f0 from f producing a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Sub(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return Fixed{fp: f.fp - f0.fp}
}

// Abs returns the absolute value of f. If f is NaN, NaN is returned
func (f Fixed) Abs() Fixed {
	if f.IsNaN() {
		return NaN
	}
	if f.Sign() >= 0 {
		return f
	}
	f0 := Fixed{fp: f.fp * -1}
	return f0
}

func abs(i int64) int64 {
	if i >= 0 {
		return i
	}
	return i * -1
}

// Mul multiplies f by f0 returning a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Mul(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}

	fp_a := f.fp / pow7
	fp_b := f.fp % pow7

	fp0_a := f0.fp / pow7
	fp0_b := f0.fp % pow7

	var result int64

	if fp0_a != 0 {
		result = fp_a*fp0_a*pow7 + fp_b*fp0_a
	}
	if fp0_b != 0 {
		result = result + (fp_a * fp0_b) + ((fp_b)*fp0_b)/pow7
	}

	return Fixed{fp: result}
}

// Div divides f by f0 returning a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Div(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return NewF(f.Float() / f0.Float())
}

// Round returns a rounded (half-up, away from zero) to n decimal places
func (f Fixed) Round(n int) Fixed {
	if f.IsNaN() {
		return NaN
	}

	round := .5
	if f.fp < 0 {
		round = -0.5
	}

	f0 := f.Frac()
	f0 = f0*math.Pow10(n) + round
	f0 = float64(int(f0)) / math.Pow10(n)

	return NewF(float64(f.Int()) + f0)
}

// Equal returns true if the f == f0. If either operand is NaN, false is returned. Use IsNaN() to test for NaN
func (f Fixed) Equal(f0 Fixed) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}
	return f.Cmp(f0) == 0
}

// GreaterThan tests Cmp() for 1
func (f Fixed) GreaterThan(f0 Fixed) bool {
	return f.Cmp(f0) == 1
}

// GreaterThaOrEqual tests Cmp() for 1 or 0
func (f Fixed) GreaterThanOrEqual(f0 Fixed) bool {
	cmp := f.Cmp(f0)
	return cmp == 1 || cmp == 0
}

// LessThan tests Cmp() for -1
func (f Fixed) LessThan(f0 Fixed) bool {
	return f.Cmp(f0) == -1
}

// LessThan tests Cmp() for -1 or 0
func (f Fixed) LessThanOrEqual(f0 Fixed) bool {
	cmp := f.Cmp(f0)
	return cmp == -1 || cmp == 0
}

// Cmp compares two Fixed. If f == f0, return 0. If f > f0, return 1. If f < f0, return -1. If both are NaN, return 0. If f is NaN, return 1. If f0 is NaN, return -1
func (f Fixed) Cmp(f0 Fixed) int {
	if f.IsNaN() && f0.IsNaN() {
		return 0
	}
	if f.IsNaN() {
		return 1
	}
	if f0.IsNaN() {
		return -1
	}

	if f.fp == f0.fp {
		return 0
	}
	if f.fp < f0.fp {
		return -1
	}
	return 1
}

// String converts a Fixed to a string, dropping trailing zeros
func (f Fixed) String() string {
	s, point := f.tostr()
	if point == -1 {
		return s
	}
	index := len(s) - 1
	for ; index != point; index-- {
		if s[index] != '0' {
			return s[:index+1]
		}
	}
	return s[:point]
}

// StringN converts a Fixed to a String with a specified number of decimal places, truncating as required
func (f Fixed) StringN(decimals int) string {

	s, point := f.tostr()

	if point == -1 {
		return s
	}
	if decimals == 0 {
		return s[:point]
	} else {
		return s[:point+decimals+1]
	}
}

func (f Fixed) tostr() (string, int) {
	fp := f.fp
	if fp == 0 {
		return "0." + zeros, 1
	}
	if fp == nan {
		return "NaN", -1
	}

	b := make([]byte, 24)
	b = itoa(b, fp)

	return string(b), len(b) - nPlaces - 1
}

func itoa(buf []byte, val int64) []byte {
	neg := val < 0
	if neg {
		val = val * -1
	}

	i := len(buf) - 1
	idec := i - nPlaces
	for val >= 10 || i >= idec {
		buf[i] = byte(val%10 + '0')
		i--
		if i == idec {
			buf[i] = '.'
			i--
		}
		val /= 10
	}
	buf[i] = byte(val + '0')
	if neg {
		i--
		buf[i] = '-'
	}
	return buf[i:]
}

// Int return the integer portion of the Fixed, or 0 if NaN
func (f Fixed) Int() int64 {
	if f.IsNaN() {
		return 0
	}
	return f.fp / pow7
}

// Frac return the fractional portion of the Fixed, or NaN if NaN
func (f Fixed) Frac() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp%pow7) / float64(pow7)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (f *Fixed) UnmarshalBinary(data []byte) error {
	fp, n := binary.Varint(data)
	if n < 0 {
		return errFormat
	}
	f.fp = fp
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (f Fixed) MarshalBinary() (data []byte, err error) {
	var buffer [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buffer[:], f.fp)
	return buffer[:n], nil
}

// WriteTo write the Fixed to an io.Writer, returning the number of bytes written
func (f Fixed) WriteTo(w io.ByteWriter) error {
	return writeVarint(w, f.fp)
}

// ReadFrom reads a Fixed from an io.Reader
func ReadFrom(r io.ByteReader) (Fixed, error) {
	fp, err := binary.ReadVarint(r)
	if err != nil {
		return NaN, err
	}
	return Fixed{fp: fp}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Fixed) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s == "null" {
		return nil
	}

	fixed, err := NewSErr(s)
	*f = fixed
	if err != nil {
		return fmt.Errorf("Error decoding string '%s': %s", s, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (f Fixed) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}
