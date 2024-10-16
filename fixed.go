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

type Places interface {
	places() int
	scale() int64
	zeros() string
	max() float64
}

const _MAX = 999999999999999999

type Places7 struct {
}
func (Places7)places() int { return 7; }
func (Places7)scale() int64 { return 10 * 10 * 10 * 10 * 10 * 10 * 10; }
func (Places7)zeros() string { return "0000000"; }
func (Places7)max() float64 { return float64(_MAX)/(10 * 10 * 10 * 10 * 10 * 10 * 10); }

var places7 = Places7{}

// Fixed is a fixed precision 38.24 number (supports 11.7 digits). It supports NaN.
type FixedN[T Places] struct {
	places T
	fp int64
}

type Fixed7 = FixedN[Places7]
type Fixed = Fixed7

// the following constants can be changed to configure a different number of decimal places - these are
// the only required changes. only 18 significant digits are supported due to NaN

const MAX = float64(99999999999.9999999)
const nan = int64(1<<63 - 1)

var NaN = Fixed{fp: nan}
var ZERO = Fixed{fp: 0}

var errTooLarge = errors.New("significand too large")
var errFormat = errors.New("invalid encoding")

// NewS creates a new Fixed from a string, returning NaN if the string could not be parsed
func NewS(s string) Fixed {
	f, _ := NewSErr(s)
	return f
}
func NewSN[T Places](s string) Fixed {
	f, _ := NewSErr(s)
	return f
}

func NewSErr(s string) (Fixed, error) {
	f, err := NewSNErr[Places7](s,places7)
	return f,err
}

func _NaN[T Places]() FixedN[T] {
	return FixedN[T]{fp:nan}
}

// NewSErr creates a new Fixed from a string, returning NaN, and error if the string could not be parsed
func NewSNErr[T Places](s string, places T) (FixedN[T], error) {
	nPlaces := places.places()

	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return _NaN[T](), err
		}
		return NewFN[T](f,places), nil
	}
	if "NaN" == s {
		return _NaN[T](), nil
	}
	period := strings.Index(s, ".")
	var i int64
	var f int64
	var sign int64 = 1
	var err error
	if period == -1 {
		i, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return _NaN[T](), errors.New("cannot parse")
		}
		if i < 0 {
			sign = -1
			i = i * -1
		}
	} else {
		if len(s[:period]) > 0 {
			i, err = strconv.ParseInt(s[:period], 10, 64)
			if err != nil {
				return _NaN[T](), errors.New("cannot parse")
			}
			if i < 0 || s[0] == '-' {
				sign = -1
				i = i * -1
			}
		}
		fs := s[period+1:]
		fs = fs + places.zeros()[:max(0, nPlaces-len(fs))]
		f, err = strconv.ParseInt(fs[0:nPlaces], 10, 64)
		if err != nil {
			return _NaN[T](), errors.New("cannot parse")
		}
	}
	if float64(i) > MAX {
		return _NaN[T](), errTooLarge
	}
	return FixedN[T]{fp: sign * (i*places.scale() + f)}, nil
}

// Parse creates a new Fixed from a string, returning NaN, and error if the string could not be parsed. Same as NewSErr
// but more standard naming
func Parse(s string) (Fixed, error) {
	return NewSErr(s)
}

func ParseN[T Places](s string,t T) (FixedN[T], error) {
	return NewSNErr(s,t)
}

// MustParse creates a new Fixed from a string, and panics if the string could not be parsed
func MustParse(s string) Fixed {
	f, err := NewSErr(s)
	if err != nil {
		panic(err)
	}
	return f
}
func MustParseN[T Places](s string,t T) FixedN[T] {
	f, err := NewSNErr(s,t)
	if err != nil {
		panic(err)
	}
	return f
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewF creates a Fixed from an float64, rounding at the 8th decimal place
func NewF(f float64) Fixed {
	return NewFN[Places7](f,places7);
}

// NewF creates a Fixed from an float64, rounding at the 8th decimal place
func NewFN[T Places](f float64,t Places) FixedN[T] {
	if math.IsNaN(f) {
		return _NaN[T]();
	}
	if f >= MAX || f <= -MAX {
		return _NaN[T]();
	}
	round := .5
	if f < 0 {
		round = -0.5
	}

	return FixedN[T]{fp: int64(f*float64(t.scale()) + round)}
}

// NewI creates a Fixed for an integer, moving the decimal point n places to the left
// For example, NewI(123,1) becomes 12.3. If n > 7, the value is truncated
func NewI(i int64, n uint) Fixed {
	return NewIN[Places7](i,n,places7)
}
	// NewI creates a Fixed for an integer, moving the decimal point n places to the left
// For example, NewI(123,1) becomes 12.3. If n > 7, the value is truncated
func NewIN[T Places](i int64, n uint, t T) Fixed {
	nPlaces := uint(t.places())
	if n > nPlaces {
		i = i / int64(math.Pow10(int(n-nPlaces)))
		n = nPlaces
	}

	i = i * int64(math.Pow10(int(nPlaces-n)))

	return Fixed{fp: i}
}

func (f FixedN[T]) IsNaN() bool {
	return f.fp == nan
}

func (f FixedN[T]) IsZero() bool {
	return f.fp == 0
}

// Sign returns:
//
//	-1 if f <  0
//	 0 if f == 0 or NaN
//	+1 if f >  0
func (f FixedN[T]) Sign() int {
	if f.IsNaN() {
		return 0
	}
	return f.Cmp(FixedN[T]{fp:0})
}

// Float converts the Fixed to a float64
func (f FixedN[T]) Float() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp) / float64(f.places.scale())
}

// Add adds f0 to f producing a Fixed. If either operand is NaN, NaN is returned
func (f FixedN[T]) Add(f0 FixedN[T]) FixedN[T] {
	if f.IsNaN() || f0.IsNaN() {
		return _NaN[T]()
	}
	return FixedN[T]{fp: f.fp + f0.fp}
}

// Sub subtracts f0 from f producing a Fixed. If either operand is NaN, NaN is returned
func (f FixedN[T]) Sub(f0 FixedN[T]) FixedN[T] {
	if f.IsNaN() || f0.IsNaN() {
		return _NaN[T]()
	}
	return FixedN[T]{fp: f.fp - f0.fp}
}

// Abs returns the absolute value of f. If f is NaN, NaN is returned
func (f FixedN[T]) Abs() FixedN[T] {
	if f.IsNaN() {
		return _NaN[T]()
	}
	if f.Sign() >= 0 {
		return f
	}
	f0 := FixedN[T]{fp: f.fp * -1}
	return f0
}

func abs(i int64) int64 {
	if i >= 0 {
		return i
	}
	return i * -1
}

// Mul multiplies f by f0 returning a Fixed. If either operand is NaN, NaN is returned
func (f FixedN[T]) Mul(f0 FixedN[T]) FixedN[T] {
	if f.IsNaN() || f0.IsNaN() {
		return _NaN[T]()
	}

	scale := f.places.scale()

	fp_a := f.fp / scale
	fp_b := f.fp % scale

	fp0_a := f0.fp / scale
	fp0_b := f0.fp % scale

	var _sign = int64(f.Sign()*f0.Sign())

	var result int64

	if fp0_a != 0 {
		result = fp_a*fp0_a*scale + fp_b*fp0_a
	}
	if fp0_b != 0 {
		result = result + (fp_a * fp0_b) + ((fp_b)*fp0_b+5*_sign*(scale/10))/scale
	}

	return FixedN[T]{fp: result}
}

// Div divides f by f0 returning a Fixed. If either operand is NaN, NaN is returned
func (f FixedN[T]) Div(f0 FixedN[T]) FixedN[T] {
	if f.IsNaN() || f0.IsNaN() {
		return _NaN[T]()
	}
	return NewFN[T](f.Float() / f0.Float(),f.places)
}

func sign(fp int64) int64 {
	if fp < 0 {
		return -1
	}
	return 1
}

// Round returns a rounded (half-up, away from zero) to n decimal places
func (f FixedN[T]) Round(n int) FixedN[T] {
	if f.IsNaN() {
		return _NaN[T]()
	}

	scale := f.places.scale()
	nPlaces := f.places.places()

	fraction := f.fp % scale
	f0 := fraction / int64(math.Pow10(nPlaces-n-1))
	digit := abs(f0 % 10)
	f0 = (f0 / 10)
	if digit >= 5 {
		f0 += 1 * sign(f.fp)
	}
	f0 = f0 * int64(math.Pow10(nPlaces-n))

	intpart := f.fp - fraction
	fp := intpart + f0

	return FixedN[T]{fp: fp}
}

// Equal returns true if the f == f0. If either operand is NaN, false is returned. Use IsNaN() to test for NaN
func (f FixedN[T]) Equal(f0 FixedN[T]) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}
	return f.Cmp(f0) == 0
}

// GreaterThan tests Cmp() for 1
func (f FixedN[T]) GreaterThan(f0 FixedN[T]) bool {
	return f.Cmp(f0) == 1
}

// GreaterThaOrEqual tests Cmp() for 1 or 0
func (f FixedN[T]) GreaterThanOrEqual(f0 FixedN[T]) bool {
	cmp := f.Cmp(f0)
	return cmp == 1 || cmp == 0
}

// LessThan tests Cmp() for -1
func (f FixedN[T]) LessThan(f0 FixedN[T]) bool {
	return f.Cmp(f0) == -1
}

// LessThan tests Cmp() for -1 or 0
func (f FixedN[T]) LessThanOrEqual(f0 FixedN[T]) bool {
	cmp := f.Cmp(f0)
	return cmp == -1 || cmp == 0
}

// Cmp compares two Fixed. If f == f0, return 0. If f > f0, return 1. If f < f0, return -1. If both are NaN, return 0. If f is NaN, return 1. If f0 is NaN, return -1
func (f FixedN[T]) Cmp(f0 FixedN[T]) int {
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
func (f FixedN[T]) String() string {
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
func (f FixedN[T]) StringN(decimals int) string {

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

func (f FixedN[T]) tostr() (string, int) {
	fp := f.fp
	if fp == 0 {
		return "0." + f.places.zeros(), 1
	}
	if fp == nan {
		return "NaN", -1
	}

	b := make([]byte, 24)
	b = itoa(b, fp, f.places.places())

	return string(b), len(b) - f.places.places() - 1
}

func itoa(buf []byte, val int64,nPlaces int) []byte {
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
func (f FixedN[T]) Int() int64 {
	if f.IsNaN() {
		return 0
	}
	return f.fp / f.places.scale()
}

// Frac return the fractional portion of the Fixed, or NaN if NaN
func (f FixedN[T]) Frac() float64 {
	scale := f.places.scale()
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp%scale) / float64(scale)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (f *FixedN[T]) UnmarshalBinary(data []byte) error {
	fp, n := binary.Varint(data)
	if n < 0 {
		return errFormat
	}
	f.fp = fp
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (f FixedN[T]) MarshalBinary() (data []byte, err error) {
	var buffer [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buffer[:], f.fp)
	return buffer[:n], nil
}

// WriteTo write the Fixed to an io.Writer, returning the number of bytes written
func (f FixedN[T]) WriteTo(w io.ByteWriter) error {
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
func (f *FixedN[T]) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s == "null" {
		return nil
	}
	if s == "\"NaN\"" {
		*f = _NaN[T]()
		return nil
	}

	fixed, err := NewSNErr(s,f.places)
	*f = fixed
	if err != nil {
		return fmt.Errorf("Error decoding string '%s': %s", s, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (f FixedN[T]) MarshalJSON() ([]byte, error) {
	if f.IsNaN() {
		return []byte("\"NaN\""), nil
	}
	buffer := make([]byte, 24)
	return itoa(buffer, f.fp,f.places.places()), nil
}
