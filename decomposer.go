//go:build !sql_scanner
// +build !sql_scanner

package fixed

import (
	"errors"
	"math/big"
)

// See https://godoc.org/github.com/golang-sql/decomposer for the decomposer.Decimal
// interface definition.

// Decompose returns the internal decimal state into parts.
// If the provided buf has sufficient capacity, buf may be returned as the coefficient with
// the value set and length set as appropriate.
func (f Fixed) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	if f.IsNaN() {
		form = 2
		return
	}
	if f.hi == 0 && f.lo == 0 {
		return
	}

	negative = f.hi < 0 || (f.hi == 0 && f.lo < 0)

	// Combine hi and lo into single coefficient
	// coefficient = hi * scale + lo (absolute values)
	absHi := f.hi
	if absHi < 0 {
		absHi = -absHi
	}
	absLo := f.lo
	if absLo < 0 {
		absLo = -absLo
	}

	bigCoef := new(big.Int).Mul(big.NewInt(absHi), big.NewInt(scale))
	bigCoef.Add(bigCoef, big.NewInt(absLo))

	coefficient = bigCoef.Bytes() // Big-endian encoding
	exponent = -18
	return
}

// Compose sets the internal decimal value from parts. If the value cannot be
// represented then an error should be returned.
func (f *Fixed) Compose(form byte, negative bool, coefficient []byte, exponent int32) (err error) {
	if f == nil {
		return errors.New("Fixed must not be nil")
	}
	switch form {
	case 1, 2: // Infinite or NaN
		f.hi = nanHi
		f.lo = nanLo
		return nil
	case 0: // Finite
		// Continue below
	default:
		return errors.New("invalid form")
	}

	// Parse coefficient from bytes
	bigCoef := new(big.Int).SetBytes(coefficient)

	// Adjust for exponent
	dividePower := int(exponent) + 18
	if dividePower != 0 {
		ct := dividePower
		if ct < 0 {
			ct = -ct
		}
		adjuster := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(ct)), nil)
		if dividePower < 0 {
			bigCoef.Div(bigCoef, adjuster)
		} else {
			bigCoef.Mul(bigCoef, adjuster)
		}
	}

	// Split into hi and lo
	bigScale := big.NewInt(scale)
	hi := new(big.Int).Div(bigCoef, bigScale)
	lo := new(big.Int).Mod(bigCoef, bigScale)

	if !hi.IsInt64() || !lo.IsInt64() {
		return errTooLarge
	}

	f.hi = hi.Int64()
	f.lo = lo.Int64()

	if negative {
		f.hi = -f.hi
		f.lo = -f.lo
	}

	return nil
}
