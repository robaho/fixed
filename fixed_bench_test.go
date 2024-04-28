package fixed

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
)

var Result Fixed

func BenchmarkAddFixed(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f0)
	}
	Result = f1
}
func BenchmarkAddDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(1)
	f1 := decimal.NewFromFloat(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f0)
	}
}
func BenchmarkAddBigInt(b *testing.B) {
	f0 := big.NewInt(1)
	f1 := big.NewInt(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f1, f0)
	}
}
func BenchmarkAddBigFloat(b *testing.B) {
	f0 := big.NewFloat(1)
	f1 := big.NewFloat(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f1, f0)
	}
}

func BenchmarkMulFixed(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Mul(f1)
	}
}
func BenchmarkMulDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Mul(f1)
	}
}
func BenchmarkMulBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var x big.Int
	for i := 0; i < b.N; i++ {
		x.Mul(f0, f1)
	}
}
func BenchmarkMulBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var x big.Float
	for i := 0; i < b.N; i++ {
		x.Mul(f0, f1)
	}
}

func BenchmarkDivFixed(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Div(f1)
	}
}
func BenchmarkDivDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Div(f1)
	}
}
func BenchmarkDivBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var x big.Int
	for i := 0; i < b.N; i++ {
		x.Div(f0, f1)
	}
}
func BenchmarkDivBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var x big.Float
	for i := 0; i < b.N; i++ {
		x.Quo(f0, f1)
	}
}

func BenchmarkCmpFixed(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}
func BenchmarkCmpDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}
func BenchmarkCmpBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}
func BenchmarkCmpBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}

func BenchmarkStringFixed(b *testing.B) {
	f0 := NewF(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}
func BenchmarkStringNFixed(b *testing.B) {
	f0 := NewF(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.StringN(5)
	}
}
func BenchmarkStringDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}
func BenchmarkStringBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}
func BenchmarkStringBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}

func BenchmarkWriteTo(b *testing.B) {
	f0 := NewF(123456789.0)

	buf := new(bytes.Buffer)

	for i := 0; i < b.N; i++ {
		f0.WriteTo(buf)
	}
}
