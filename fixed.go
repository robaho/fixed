package fixed

// release under the terms of file license.txt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Fixed is a fixed precision 18.18 number (18 integer digits, 18 decimal digits). It supports NaN.
type Fixed struct {
	hi int64 // Integer part: -999999999999999999 to 999999999999999999
	lo int64 // Decimal part scaled by 10^18, carries same sign as hi
}

// the following constants can be changed to configure a different number of decimal places - these are
// the only required changes. only 18 significant digits are supported due to NaN

const nPlaces = 18
const scale = int64(1000000000000000000) // 10^18
const zeros = "000000000000000000"
const MAX = float64(999999999999999999.999999999999999999)
const maxHi = int64(999999999999999999) // Maximum valid hi value (18 digits)

// NaN representation: both fields set to int64 max
const nanHi = int64(1<<63 - 1)
const nanLo = int64(1<<63 - 1)

var NaN = Fixed{hi: nanHi, lo: nanLo}
var ZERO = Fixed{hi: 0, lo: 0}

var errTooLarge = errors.New("significand too large")
var errFormat = errors.New("invalid encoding")

var pow10table = [19]int64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

var errPow10Overflow = errors.New("pow10: exponent out of int64 range")

func ipow10(n int) (int64, error) {
	if n < 0 {
		return 0, errPow10Overflow
	}
	if n > 18 {
		return pow10table[18], errPow10Overflow
	}
	return pow10table[n], nil
}

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
	var hi int64
	var lo int64
	var sign int64 = 1
	var err error
	if period == -1 {
		hi, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return NaN, errors.New("cannot parse")
		}
		if hi < 0 {
			sign = -1
			hi = hi * -1
		}
	} else {
		intPart := s[:period]
		if len(intPart) > 0 && intPart != "-" && intPart != "+" {
			hi, err = strconv.ParseInt(intPart, 10, 64)
			if err != nil {
				return NaN, errors.New("cannot parse")
			}
			if hi < 0 {
				sign = -1
				hi = hi * -1
			}
		}
		// Check for sign prefix even if no integer digits
		if len(intPart) > 0 && intPart[0] == '-' {
			sign = -1
		}
		fs := s[period+1:]
		fs = fs + zeros[:max(0, nPlaces-len(fs))]
		lo, err = strconv.ParseInt(fs[0:nPlaces], 10, 64)
		if err != nil {
			return NaN, errors.New("cannot parse")
		}
	}
	// Check for overflow - hi must fit within 18 digits
	if hi > maxHi {
		return NaN, errTooLarge
	}
	return Fixed{hi: sign * hi, lo: sign * lo}, nil
}

// Parse creates a new Fixed from a string, returning NaN, and error if the string could not be parsed. Same as NewSErr
// but more standard naming
func Parse(s string) (Fixed, error) {
	return NewSErr(s)
}

// MustParse creates a new Fixed from a string, and panics if the string could not be parsed
func MustParse(s string) Fixed {
	f, err := NewSErr(s)
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

// normalize ensures lo is within proper range and adjusts hi accordingly
func normalize(hi, lo int64) (int64, int64) {
	if lo >= scale {
		hi += lo / scale
		lo = lo % scale
	} else if lo <= -scale {
		hi += lo / scale // lo/scale is negative
		lo = lo % scale
	}

	// Handle sign consistency: both parts should have same sign
	// Special case: hi=0 with non-zero lo keeps lo's sign
	if hi > 0 && lo < 0 {
		hi--
		lo += scale
	} else if hi < 0 && lo > 0 {
		hi++
		lo -= scale
	}

	return hi, lo
}

// NewF creates a Fixed from a float64, zero-allocation.
func NewF(f float64) Fixed {
	if math.IsNaN(f) {
		return NaN
	}
	if f >= MAX || f <= -MAX {
		return NaN
	}
	if f == 0 {
		return ZERO
	}

	// Format into stack buffer (no heap allocation)
	var buf [64]byte
	b := strconv.AppendFloat(buf[:0], f, 'f', -1, 64)

	// Parse bytes directly
	i := 0
	neg := false
	if b[0] == '-' {
		neg = true
		i++
	}

	// Scan integer digits
	var hi int64
	for i < len(b) && b[i] != '.' {
		hi = hi*10 + int64(b[i]-'0')
		i++
	}

	// Scan fractional digits
	var lo int64
	var nFrac int
	if i < len(b) && b[i] == '.' {
		i++ // skip '.'
		for i < len(b) && nFrac < nPlaces {
			lo = lo*10 + int64(b[i]-'0')
			nFrac++
			i++
		}
	}

	// Pad lo to nPlaces digits
	for nFrac < nPlaces {
		lo *= 10
		nFrac++
	}

	if hi > maxHi {
		return NaN
	}

	if neg {
		hi = -hi
		lo = -lo
	}

	return Fixed{hi: hi, lo: lo}
}

// NewFSlow creates a Fixed from a float64 via string conversion (allocates).
func NewFSlow(f float64) Fixed {
	if math.IsNaN(f) {
		return NaN
	}
	if f >= MAX || f <= -MAX {
		return NaN
	}

	// Convert to string with 15 significant figures (float64's precision limit)
	// then parse the string to avoid float64 precision artifacts
	s := strconv.FormatFloat(f, 'f', -1, 64)

	// Parse and return - this handles the conversion cleanly
	result, err := NewSErr(s)
	if err != nil {
		// Fallback: direct conversion
		hi := int64(math.Trunc(f))
		frac := f - float64(hi)
		lo := int64(frac * float64(scale))
		hi, lo = normalize(hi, lo)
		return Fixed{hi: hi, lo: lo}
	}

	return result
}

// NewI creates a Fixed for an integer, moving the decimal point n places to the left
// For example, NewI(123,1) becomes 12.3. If n > 18, the value is truncated
func NewI(i int64, n uint) Fixed {
	if n > nPlaces {
		p, err := ipow10(int(n - nPlaces))
		if err != nil {
			panic(err)
		}
		i = i / p
		n = nPlaces
	}

	// Split i into integer and decimal portions based on n
	divisor, _ := ipow10(int(n))
	hi := i / divisor
	remainder := i % divisor

	// Scale decimal portion to 18 digits
	p, _ := ipow10(int(nPlaces - n))
	lo := remainder * p

	// Check for overflow
	if hi > maxHi || hi < -maxHi {
		return NaN
	}

	return Fixed{hi: hi, lo: lo}
}

func (f Fixed) IsNaN() bool {
	return f.hi == nanHi && f.lo == nanLo
}

func (f Fixed) IsZero() bool {
	return f.Equal(ZERO)
}

// Sign returns:
//
//	-1 if f <  0
//	 0 if f == 0 or NaN
//	+1 if f >  0
func (f Fixed) Sign() int {
	if f.IsNaN() {
		return 0
	}
	if f.hi < 0 || (f.hi == 0 && f.lo < 0) {
		return -1
	}
	if f.hi > 0 || (f.hi == 0 && f.lo > 0) {
		return 1
	}
	return 0
}

// Float converts the Fixed to a float64
func (f Fixed) Float() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.hi) + float64(f.lo)/float64(scale)
}

// Add adds f0 to f producing a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Add(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	hi := f.hi + f0.hi
	lo := f.lo + f0.lo
	hi, lo = normalize(hi, lo)

	// Check for overflow
	if hi > maxHi || hi < -maxHi {
		return NaN
	}

	return Fixed{hi: hi, lo: lo}
}

// Sub subtracts f0 from f producing a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Sub(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	hi := f.hi - f0.hi
	lo := f.lo - f0.lo
	hi, lo = normalize(hi, lo)

	// Check for overflow
	if hi > maxHi || hi < -maxHi {
		return NaN
	}

	return Fixed{hi: hi, lo: lo}
}

// Abs returns the absolute value of f. If f is NaN, NaN is returned
func (f Fixed) Abs() Fixed {
	if f.IsNaN() {
		return NaN
	}
	if f.Sign() >= 0 {
		return f
	}
	// Negate both parts
	return Fixed{hi: -f.hi, lo: -f.lo}
}

func abs(i int64) int64 {
	if i >= 0 {
		return i
	}
	return i * -1
}

// half is 10^9, used as the base for schoolbook multiplication in Mul.
const half = int64(1000000000)

// Mul multiplies f by f0 returning a Fixed. If either operand is NaN, NaN is returned.
// Zero-allocation implementation using base-10^9 schoolbook multiplication.
func (f Fixed) Mul(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}

	// Determine result sign, then work with absolute values
	signA := f.Sign()
	signB := f0.Sign()
	if signA == 0 || signB == 0 {
		return ZERO
	}
	negative := signA != signB

	aHi := f.hi
	aLo := f.lo
	if aHi < 0 {
		aHi = -aHi
	}
	if aLo < 0 {
		aLo = -aLo
	}

	bHi := f0.hi
	bLo := f0.lo
	if bHi < 0 {
		bHi = -bHi
	}
	if bLo < 0 {
		bLo = -bLo
	}

	// Decompose each operand into 4 digits in base half (10^9):
	//   value = (d[0]*half + d[1]) * scale + (d[2]*half + d[3])
	// where scale = half * half = 10^18
	var a [4]int64
	a[0] = aHi / half
	a[1] = aHi % half
	a[2] = aLo / half
	a[3] = aLo % half

	var b [4]int64
	b[0] = bHi / half
	b[1] = bHi % half
	b[2] = bLo / half
	b[3] = bLo % half

	// Convolution: p[k] = sum of a[i]*b[j] where i+j=k
	// The full product has digits p[0]..p[6]. Dividing by scale (= half^2)
	// shifts by 2 digit positions, so the result digits are p[1]..p[4].
	// p[0] must be zero (else overflow). p[5],p[6] are truncated.
	// Each product a[i]*b[j] < (10^9)^2 = 10^18, at most 4 terms per p[k],
	// so p[k] < 4*10^18 < math.MaxInt64. No overflow in int64.
	var p [7]int64
	p[0] = a[0] * b[0]
	p[1] = a[0]*b[1] + a[1]*b[0]
	p[2] = a[0]*b[2] + a[1]*b[1] + a[2]*b[0]
	p[3] = a[0]*b[3] + a[1]*b[2] + a[2]*b[1] + a[3]*b[0]
	p[4] = a[1]*b[3] + a[2]*b[2] + a[3]*b[1]
	p[5] = a[2]*b[3] + a[3]*b[2]
	p[6] = a[3] * b[3]

	// Carry-propagate from p[6] up to p[0]
	p[5] += p[6] / half
	p[6] = p[6] % half
	p[4] += p[5] / half
	p[5] = p[5] % half
	p[3] += p[4] / half
	p[4] = p[4] % half
	p[2] += p[3] / half
	p[3] = p[3] % half
	p[1] += p[2] / half
	p[2] = p[2] % half
	p[0] += p[1] / half
	p[1] = p[1] % half

	// Overflow check: the top digit must be zero after carry
	if p[0] != 0 {
		return NaN
	}

	// Reconstruct hi and lo from result digits [1..4]
	// The product has 7 digits (p[0]..p[6]). Dividing by scale=half^2 shifts
	// by 2 digit positions, so integer part = p[0]*B^4 + p[1]*B^3 + p[2]*B^2 + p[3]*B + p[4].
	// With p[0]=0: hi = p[1]*half + p[2], lo = p[3]*half + p[4].
	// Digits p[5], p[6] are truncated (division toward zero).
	resHi := p[1]*half + p[2]
	resLo := p[3]*half + p[4]

	// For negative results with a non-zero truncated remainder, adjust to match
	// floor division (round toward -infinity) semantics used by MulSlow's big.Int.Div.
	// Truncation gives |result|, floor division gives |result|+1 when remainder > 0.
	if negative && (p[5] != 0 || p[6] != 0) {
		resLo++
		if resLo >= scale {
			resLo -= scale
			resHi++
		}
	}

	// Check if result exceeds valid range
	if resHi > maxHi {
		return NaN
	}

	// Apply sign
	if negative && (resHi != 0 || resLo != 0) {
		resHi = -resHi
		resLo = -resLo
	}

	return Fixed{hi: resHi, lo: resLo}
}

// MulSlow multiplies f by f0 returning a Fixed. If either operand is NaN, NaN is returned
// Uses math/big for all multiplications — correct but allocates.
func (f Fixed) MulSlow(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}

	// Use math/big for all multiplications to avoid overflow
	// Convert to single big integer (value * scale), multiply, then split back
	bigScale := big.NewInt(scale)

	// value1 = hi1 * scale + lo1
	value1 := new(big.Int).Mul(big.NewInt(f.hi), bigScale)
	value1.Add(value1, big.NewInt(f.lo))

	// value2 = hi2 * scale + lo2
	value2 := new(big.Int).Mul(big.NewInt(f0.hi), bigScale)
	value2.Add(value2, big.NewInt(f0.lo))

	// product = value1 * value2 / scale (to maintain scale)
	product := new(big.Int).Mul(value1, value2)
	product.Div(product, bigScale)

	// Split back into hi and lo
	hi := new(big.Int).Div(product, bigScale)
	lo := new(big.Int).Mod(product, bigScale)

	// Check overflow
	if !hi.IsInt64() || !lo.IsInt64() {
		return NaN
	}

	// Normalize to ensure sign consistency
	resHi, resLo := normalize(hi.Int64(), lo.Int64())

	// Check if result exceeds valid range
	if resHi > maxHi || resHi < -maxHi {
		return NaN
	}

	return Fixed{hi: resHi, lo: resLo}
}

// Div divides f by f0 returning a Fixed. If either operand is NaN, NaN is returned.
// Zero-allocation implementation using base-10^9 long division (Knuth's Algorithm D).
func (f Fixed) Div(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	if f0.hi == 0 && f0.lo == 0 {
		return NaN // division by zero
	}

	// Determine result sign, then work with absolute values
	signA := f.Sign()
	signB := f0.Sign()
	if signA == 0 {
		return ZERO
	}
	negative := signA != signB

	aHi := f.hi
	aLo := f.lo
	if aHi < 0 {
		aHi = -aHi
	}
	if aLo < 0 {
		aLo = -aLo
	}

	bHi := f0.hi
	bLo := f0.lo
	if bHi < 0 {
		bHi = -bHi
	}
	if bLo < 0 {
		bLo = -bLo
	}

	// Decompose dividend into base-B (B=half=10^9) digits.
	// Dividend = (aHi*scale + aLo) * scale = aHi*B^4 + aLo*B^2
	// 7 digits: u[0]=0 (leading zero for Algorithm D), u[1..4] from value, u[5..6]=0 from ×scale
	var u [7]int64
	u[1] = aHi / half
	u[2] = aHi % half
	u[3] = aLo / half
	u[4] = aLo % half

	// Decompose divisor: bHi*scale + bLo = bHi*B^2 + bLo
	// Up to 4 digits, strip leading zeros to get n digits
	var vBuf [4]int64
	vBuf[0] = bHi / half
	vBuf[1] = bHi % half
	vBuf[2] = bLo / half
	vBuf[3] = bLo % half

	vStart := 0
	for vStart < 3 && vBuf[vStart] == 0 {
		vStart++
	}
	n := 4 - vStart

	var v [4]int64
	for i := 0; i < n; i++ {
		v[i] = vBuf[vStart+i]
	}

	const m = 6 // dividend has m+1 = 7 digits (u[0]..u[6])

	var q [7]int64 // quotient digits

	if n == 1 {
		// Simple single-digit long division
		rem := int64(0)
		for i := 0; i <= m; i++ {
			cur := rem*half + u[i]
			q[i] = cur / v[0]
			rem = cur % v[0]
		}

		// Overflow check: first 3 quotient digits must be zero
		if q[0] != 0 || q[1] != 0 || q[2] != 0 {
			return NaN
		}

		resHi := q[3]*half + q[4]
		resLo := q[5]*half + q[6]

		// Round to nearest: if 2*rem >= divisor, round up (away from zero)
		if 2*rem >= v[0] {
			resLo++
			if resLo >= scale {
				resLo -= scale
				resHi++
			}
		}

		if resHi > maxHi {
			return NaN
		}

		if negative && (resHi != 0 || resLo != 0) {
			resHi = -resHi
			resLo = -resLo
		}

		return Fixed{hi: resHi, lo: resLo}
	}

	// n >= 2: Knuth's Algorithm D

	// D1. Normalize: multiply u and v by d = B/(v[0]+1) so that v[0] >= B/2.
	// All intermediate products fit in int64: digit*d < B*B = 10^18 < maxInt64.
	d := half / (v[0] + 1)

	// Multiply u[0..m] by d (right to left, propagating carry)
	carry := int64(0)
	for i := m; i >= 0; i-- {
		tmp := u[i]*d + carry
		u[i] = tmp % half
		carry = tmp / half
	}
	// carry == 0: u[0] was 0 and absorbs any carry from u[1..6]

	// Multiply v[0..n-1] by d
	carry = 0
	for i := n - 1; i >= 0; i-- {
		tmp := v[i]*d + carry
		v[i] = tmp % half
		carry = tmp / half
	}
	// After normalization: v[0] >= floor(B/2)

	// D2-D7. Main loop: compute quotient digits q[0..m-n]
	for j := 0; j <= m-n; j++ {
		// D3. Calculate trial quotient q̂
		qhat := (u[j]*half + u[j+1]) / v[0]
		rhat := (u[j]*half + u[j+1]) % v[0]

		// Refine q̂ using second divisor digit
		for qhat >= half || qhat*v[1] > rhat*half+u[j+2] {
			qhat--
			rhat += v[0]
			if rhat >= half {
				break
			}
		}

		// D4. Multiply and subtract: u[j..j+n] -= qhat * v[0..n-1]
		carry = 0
		borrow := int64(0)
		for k := n - 1; k >= 0; k-- {
			p := qhat*v[k] + carry
			carry = p / half
			pLo := p % half

			diff := u[j+1+k] - pLo - borrow
			if diff < 0 {
				diff += half
				borrow = 1
			} else {
				borrow = 0
			}
			u[j+1+k] = diff
		}
		u[j] -= carry + borrow

		// D5. Set quotient digit
		q[j] = qhat

		if u[j] < 0 {
			// D6. Add back (rare: qhat was one too large)
			q[j]--
			carry = 0
			for k := n - 1; k >= 0; k-- {
				sum := u[j+1+k] + v[k] + carry
				u[j+1+k] = sum % half
				carry = sum / half
			}
			u[j] += carry
		}
	}

	// Quotient is in q[0..m-n], total qLen = m-n+1 = 7-n digits.
	// For the result to fit in 4 base-B digits, leading digits must be zero.
	qLen := m - n + 1
	for i := 0; i < qLen-4; i++ {
		if q[i] != 0 {
			return NaN
		}
	}

	// Extract 4-digit quotient (pad with leading zeros if qLen < 4)
	var qd [4]int64
	for i := 0; i < 4; i++ {
		idx := qLen - 4 + i
		if idx >= 0 {
			qd[i] = q[idx]
		}
	}

	resHi := qd[0]*half + qd[1]
	resLo := qd[2]*half + qd[3]

	// Round to nearest: compare 2×remainder with divisor.
	// Both are still multiplied by normalization factor d, so comparison is valid.
	// Remainder is in u[m-n+1..m] (n digits), divisor is v[0..n-1] (n digits).
	roundUp := false
	carry2 := int64(0)
	var rem2 [4]int64
	for i := n - 1; i >= 0; i-- {
		tmp := u[m-n+1+i]*2 + carry2
		rem2[i] = tmp % half
		carry2 = tmp / half
	}
	if carry2 > 0 {
		roundUp = true
	} else {
		for i := 0; i < n; i++ {
			if rem2[i] > v[i] {
				roundUp = true
				break
			} else if rem2[i] < v[i] {
				break
			}
		}
	}

	if roundUp {
		resLo++
		if resLo >= scale {
			resLo -= scale
			resHi++
		}
	}

	if resHi > maxHi {
		return NaN
	}

	if negative && (resHi != 0 || resLo != 0) {
		resHi = -resHi
		resLo = -resLo
	}

	return Fixed{hi: resHi, lo: resLo}
}

// DivSlow divides f by f0 returning a Fixed. If either operand is NaN, NaN is returned
// Uses arbitrary precision math/big for accurate division
func (f Fixed) DivSlow(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	if f0.hi == 0 && f0.lo == 0 {
		return NaN // division by zero
	}

	// Convert to big.Int for full precision
	dividend := new(big.Int).Mul(big.NewInt(f.hi), big.NewInt(scale))
	dividend.Add(dividend, big.NewInt(f.lo))

	divisor := new(big.Int).Mul(big.NewInt(f0.hi), big.NewInt(scale))
	divisor.Add(divisor, big.NewInt(f0.lo))

	// Scale dividend by 10^18 for proper decimal places
	dividend.Mul(dividend, big.NewInt(scale))

	// Perform division with rounding toward nearest
	// Use QuoRem for truncated division (toward zero) instead of DivMod (Euclidean)
	result := new(big.Int)
	remainder := new(big.Int)
	result.QuoRem(dividend, divisor, remainder)

	// Round to nearest: if abs(2*remainder) >= abs(divisor), round away from zero
	absRemainder := new(big.Int).Abs(remainder)
	absRemainder.Mul(absRemainder, big.NewInt(2))
	absDivisor := new(big.Int).Abs(divisor)

	if absRemainder.Cmp(absDivisor) >= 0 {
		if result.Sign() >= 0 {
			result.Add(result, big.NewInt(1))
		} else {
			result.Sub(result, big.NewInt(1))
		}
	}

	// Split result back into hi and lo using truncated division
	// to avoid Euclidean division sign issues
	bigScale := big.NewInt(scale)
	hi := new(big.Int).Div(result, bigScale)

	// Calculate lo as: lo = result - hi*scale (preserves sign correctly)
	lo := new(big.Int).Mul(hi, bigScale)
	lo.Sub(result, lo)

	if !hi.IsInt64() || !lo.IsInt64() {
		return NaN
	}

	// Normalize to ensure sign consistency
	resHi, resLo := normalize(hi.Int64(), lo.Int64())
	return Fixed{hi: resHi, lo: resLo}
}

func sign(fp int64) int64 {
	if fp < 0 {
		return -1
	}
	return 1
}

// Round returns a rounded (half-up, away from zero) to n decimal places
func (f Fixed) Round(n int) Fixed {
	if f.IsNaN() {
		return NaN
	}
	if n >= 18 {
		return f
	}

	if n >= 0 {
		// Rounding decimal part (lo)
		divisor, _ := ipow10(18 - n)
		remainder := f.lo % divisor
		absRemainder := remainder
		if absRemainder < 0 {
			absRemainder = -absRemainder
		}

		newLo := (f.lo / divisor) * divisor

		// Half-up rounding
		halfDivisor := divisor / 2
		if absRemainder >= halfDivisor {
			if f.lo >= 0 {
				newLo += divisor
			} else {
				newLo -= divisor
			}
		}

		hi, lo := normalize(f.hi, newLo)
		return Fixed{hi: hi, lo: lo}
	} else {
		// Rounding integer part (hi), zero out lo
		divisor, err := ipow10(-n)
		if err != nil {
			panic(err)
		}
		remainder := f.hi % divisor
		absRemainder := remainder
		if absRemainder < 0 {
			absRemainder = -absRemainder
		}

		newHi := (f.hi / divisor) * divisor

		if absRemainder >= divisor/2 {
			if f.hi >= 0 {
				newHi += divisor
			} else {
				newHi -= divisor
			}
		}

		return Fixed{hi: newHi, lo: 0}
	}
}

// Ceil returns f rounded up to n decimal places
func (f Fixed) Ceil(n int) Fixed {
	if f.IsNaN() {
		return NaN
	}
	f0 := f.Round(n)
	if f0.Cmp(f) >= 0 {
		return f0
	}
	adj := int64(1)
	if n < 0 {
		p, err := ipow10(-n)
		if err != nil {
			panic(err)
		}
		adj = adj * p
		n = 0
	}
	return f0.Add(NewI(adj, uint(n)))
}

// Floor returns f rounded down to n decimal places
func (f Fixed) Floor(n int) Fixed {
	if f.IsNaN() {
		return NaN
	}
	f0 := f.Round(n)
	if f0.Cmp(f) <= 0 {
		return f0
	}
	adj := int64(-1)
	if n < 0 {
		p, err := ipow10(-n)
		if err != nil {
			panic(err)
		}
		adj = adj * p
		n = 0
	}
	return f0.Add(NewI(adj, uint(n)))
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

	if f.hi < f0.hi {
		return -1
	}
	if f.hi > f0.hi {
		return 1
	}

	// hi parts equal, compare lo parts
	if f.lo < f0.lo {
		return -1
	}
	if f.lo > f0.lo {
		return 1
	}

	return 0
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
	if f.hi == 0 && f.lo == 0 {
		return "0." + zeros, 1
	}
	if f.IsNaN() {
		return "NaN", -1
	}

	// Determine sign
	negative := f.hi < 0 || (f.hi == 0 && f.lo < 0)

	// Work with absolute values
	absHi := f.hi
	if absHi < 0 {
		absHi = -absHi
	}
	absLo := f.lo
	if absLo < 0 {
		absLo = -absLo
	}

	// Build string right-to-left in a stack-allocated buffer.
	// Max length: '-' (1) + 18 integer digits + '.' (1) + 18 fractional digits = 38
	var buf [40]byte
	pos := len(buf)

	// Write lo digits right-to-left (always nPlaces digits, zero-padded)
	for i := 0; i < nPlaces; i++ {
		pos--
		buf[pos] = byte(absLo%10) + '0'
		absLo /= 10
	}

	// Decimal point
	pos--
	buf[pos] = '.'
	point := pos

	// Write hi digits right-to-left
	if absHi == 0 {
		pos--
		buf[pos] = '0'
	} else {
		for absHi > 0 {
			pos--
			buf[pos] = byte(absHi%10) + '0'
			absHi /= 10
		}
	}

	// Sign
	if negative {
		pos--
		buf[pos] = '-'
	}

	return string(buf[pos:]), point - pos
}

func itoa(buf []byte, val int64) []byte {
	// Note: This function is deprecated in favor of tostr()
	// but kept for backward compatibility with MarshalJSON
	// We'll use a Fixed value to leverage tostr()
	f := Fixed{hi: val / scale, lo: val % scale}
	s, _ := f.tostr()
	return []byte(s)
}

// Int return the integer portion of the Fixed, or 0 if NaN
func (f Fixed) Int() int64 {
	if f.IsNaN() {
		return 0
	}
	return f.hi
}

// Frac return the fractional portion of the Fixed, or NaN if NaN
func (f Fixed) Frac() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.lo) / float64(scale)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (f *Fixed) UnmarshalBinary(data []byte) error {
	hi, n := binary.Varint(data)
	if n <= 0 {
		return errFormat
	}

	lo, m := binary.Varint(data[n:])
	if m <= 0 {
		return errFormat
	}

	f.hi = hi
	f.lo = lo
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (f Fixed) MarshalBinary() (data []byte, err error) {
	var buffer [2 * binary.MaxVarintLen64]byte
	n := binary.PutVarint(buffer[:], f.hi)
	n += binary.PutVarint(buffer[n:], f.lo)
	return buffer[:n], nil
}

// WriteTo write the Fixed to an io.Writer, returning the number of bytes written
func (f Fixed) WriteTo(w io.ByteWriter) error {
	if err := writeVarint(w, f.hi); err != nil {
		return err
	}
	return writeVarint(w, f.lo)
}

// ReadFrom reads a Fixed from an io.Reader
func ReadFrom(r io.ByteReader) (Fixed, error) {
	hi, err := binary.ReadVarint(r)
	if err != nil {
		return NaN, err
	}

	lo, err := binary.ReadVarint(r)
	if err != nil {
		return NaN, err
	}

	return Fixed{hi: hi, lo: lo}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Fixed) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s == "null" {
		return nil
	}
	if s == "\"NaN\"" {
		*f = NaN
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
	if f.IsNaN() {
		return []byte("\"NaN\""), nil
	}
	return []byte(f.String()), nil
}
