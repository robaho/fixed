//go:build !sql_scanner
// +build !sql_scanner

package fixed

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// See https://godoc.org/github.com/golang-sql/decomposer for the decomposer.Decimal
// interface definition.

// Decompose returns the internal decimal state into parts.
// If the provided buf has sufficient capacity, buf may be returned as the coefficient with
// the value set and length set as appropriate.
func (f FixedN[T]) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	if f.fp == nan {
		form = 2
		return
	}
	if f.fp == 0 {
		return
	}
	c := f.fp
	if c < 0 {
		negative = true
		c = -c
	}
	if cap(buf) >= 8 {
		coefficient = buf[:8]
	} else {
		coefficient = make([]byte, 8)
	}
	binary.BigEndian.PutUint64(coefficient, uint64(c))
	exponent = -(int32(f.places.places()))
	return
}

// Compose sets the internal decimal value from parts. If the value cannot be
// represented then an error should be returned.
func (f *FixedN[T]) Compose(form byte, negative bool, coefficient []byte, exponent int32) (err error) {
	if f == nil {
		return errors.New("Fixed must not be nil")
	}
	switch form {
	default:
		return errors.New("invalid form")
	case 0:
		// Finite form, see below.
	case 1:
		// Infinite form, turn into NaN.
		f.fp = nan
		return nil
	case 2:
		f.fp = nan
		return nil
	}
	// Finite form.

	var c uint64
	maxi := len(coefficient) - 1
	for i := range coefficient {
		v := coefficient[maxi-i]
		if i < 8 {
			c |= uint64(v) << (uint(i) * 8)
		} else if v != 0 {
			return fmt.Errorf("coefficent too large")
		}
	}

	dividePower := int(exponent) + f.places.places()
	ct := dividePower
	if ct < 0 {
		ct = -ct
	}
	var power uint64 = 1
	for i := 0; i < ct; i++ {
		power *= 10
	}
	checkC := c
	if dividePower < 0 {
		c = c / power
		if c*power != checkC {
			return fmt.Errorf("unable to store decimal, greater then 7 decimals")
		}
	} else if dividePower > 0 {
		c = c * power
		if c/power != checkC {
			return fmt.Errorf("enable to store decimal, too large")
		}
	}
	f.fp = int64(c)
	if negative {
		f.fp = -f.fp
	}
	return nil
}
