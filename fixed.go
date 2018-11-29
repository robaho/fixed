package fixed

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// Fixed is a fixed precision 38.24 number (supports 12.7 digits). It supports NaN.
type Fixed struct {
	fp int64
}

const nPlaces = 7
const pow7 = int64(10 * 10 * 10 * 10 * 10 * 10 * 10)
const zeros = "0000000"
const nan = int64(1 << 12)

var NaN = Fixed{fp: nan}
var ZERO = Fixed{fp: 0}

const MAX = float64(99999999999.9999999)

var errTooLarge = errors.New("significand too large")

// NewS creates a new Fixed from a string, returning NaN if the string could not be parsed
func NewS(s string) Fixed {
	f, _ := NewSErr(s)
	return f
}

// NewS creates a new Fixed from a string, returning NaN, and error if the string could not be parsed
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
		fs = fs + zeros[:nPlaces-len(fs)]
		f, _ = strconv.ParseInt(fs[0:nPlaces], 10, 64)
	}
	if i > 99999999999 {
		return NaN, errTooLarge
	}
	return Fixed{fp: sign * (i*pow7 + f)}, nil
}

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

func (f Fixed) Float() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp) / float64(pow7)
}

func (f Fixed) Add(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return Fixed{fp: f.fp + f0.fp}
}

func (f Fixed) Sub(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return Fixed{fp: f.fp - f0.fp}
}

func (f Fixed) Mul(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return NewF(f.Float() * f0.Float())
}

func (f Fixed) Div(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return NewF(f.Float() / f0.Float())
}

// Round returns a rounded (half-up) to n decimal places
func (f Fixed) Round(n int) Fixed {
	if f.IsNaN() {
		return NaN
	}
	f0 := f.Frac()
	f0 = f0*math.Pow10(n) + .5
	f0 = float64(int(f0)) / math.Pow10(n)

	return NewF(float64(f.Int()) + f0)
}

func (f Fixed) Equal(f0 Fixed) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}
	return f.Cmp(f0) == 0
}

func (f Fixed) GreaterThan(f0 Fixed) bool {
	return f.Cmp(f0) == 1
}
func (f Fixed) GreaterThanOrEqual(f0 Fixed) bool {
	cmp := f.Cmp(f0)
	return cmp == 1 || cmp == 0
}
func (f Fixed) LessThan(f0 Fixed) bool {
	return f.Cmp(f0) == -1
}
func (f Fixed) LessThanOrEqual(f0 Fixed) bool {
	cmp := f.Cmp(f0)
	return cmp == -1 || cmp == 0
}

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

	fp := f.fp
	if fp == 0 {
		return "0"
	}
	if fp == nan {
		return "NaN"
	}

	var buffer [32]byte
	b := buffer[:0]

	if fp < 0 {
		b = append(b, byte('-'))
		fp *= -1
	}

	dec := fp / pow7
	frac := fp % pow7

	b = strconv.AppendInt(b, dec, 10)
	if frac == 0 {
		return string(b)
	} else {
		var buffer [32]byte
		b0 := buffer[:0]

		b = append(b, byte('.'))
		b0 = strconv.AppendInt(b0, frac, 10)
		b = append(b, []byte(zeros[:nPlaces-len(b0)])...)
		b = append(b, b0...)

		for l := len(b); l >= 0; l-- {
			if b[l-1] != '0' {
				return string(b[:l])
			}
		}
		return string(b)
	}
}

// StringN converts a Fixed to a String with a specified number of decimal places, truncating as required
func (f Fixed) StringN(decimals int) string {

	fp := f.fp
	if fp == 0 {
		if decimals == 0 {
			return "0"
		}
		return "0." + zeros[nPlaces-decimals:]
	}
	if fp == nan {
		return "NaN"
	}

	var buffer [32]byte
	b := buffer[:0]

	if fp < 0 {
		b = append(b, byte('-'))
		fp *= -1
	}

	dec := fp / pow7
	frac := fp % pow7

	b = strconv.AppendInt(b, dec, 10)
	if frac == 0 || decimals == 0 {
		return string(b)
	} else {
		var buffer [32]byte
		b0 := buffer[:0]

		b = append(b, byte('.'))
		b0 = strconv.AppendInt(b0, frac, 10)
		b = append(b, []byte(zeros[:nPlaces-len(b0)])...)
		b = append(b, b0[:decimals]...)
		return string(b)
	}
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
// do not use will be removed
func (f Fixed) ToRaw() int64 {
	return f.fp
}
// do not use will be removed
func FromRaw(i int64) Fixed {
	f := Fixed{fp: i}
	return f
}
