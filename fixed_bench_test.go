package fixed

import (
	"bytes"
	"github.com/shopspring/decimal"
	"math/big"
	"testing"
)

func BenchmarkAddFixed(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f0)
	}
	b.Log(f1)
}
func BenchmarkAddDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(1)
	f1 := decimal.NewFromFloat(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f0)
	}
	b.Log(f1)
}
func BenchmarkAddBigInt(b *testing.B) {
	f0 := big.NewInt(1)
	f1 := big.NewInt(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f1, f0)
	}
	b.Log(f1)
}
func BenchmarkAddBigFloat(b *testing.B) {
	f0 := big.NewFloat(1)
	f1 := big.NewFloat(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f1, f0)
	}
	b.Log(f1)
}

func BenchmarkMulFixed(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Mul(f1)
	}
	b.Log(f0)
}
func BenchmarkMulDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Mul(f1)
	}
	b.Log(f0)
}
func BenchmarkMulBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var x big.Int
	for i := 0; i < b.N; i++ {
		x.Mul(f0, f1)
	}
	b.Log(x)
}
func BenchmarkMulBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var x big.Float
	for i := 0; i < b.N; i++ {
		x.Mul(f0, f1)
	}
	b.Log(x)
}

func BenchmarkDivFixed(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Div(f1)
	}
	b.Log(f0)
}
func BenchmarkDivDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Div(f1)
	}
	b.Log(f0)
}
func BenchmarkDivBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var x big.Int
	for i := 0; i < b.N; i++ {
		x.Div(f0, f1)
	}
	b.Log(x)
}
func BenchmarkDivBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var x big.Float
	for i := 0; i < b.N; i++ {
		x.Quo(f0, f1)
	}
	b.Log(x)
}

func BenchmarkCmpFixed(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	var i int64
	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
	b.Log(i)
}
func BenchmarkCmpDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	var i int64
	for i := 0; i < b.N; i++ {
		i += f0.Cmp(f1)
	}
	b.Log(i)
}
func BenchmarkCmpBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var i int64
	for i := 0; i < b.N; i++ {
		i += f0.Cmp(f1)
	}
	b.Log(i)
}
func BenchmarkCmpBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var i int64
	for i := 0; i < b.N; i++ {
		i += f0.Cmp(f1)
	}
	b.Log(i)
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
	b.Log(buf.String())
}
