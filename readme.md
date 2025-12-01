**Summary**

A fixed place numeric library designed for performance.

The C++ version is available [here](https://github.com/robaho/cpp_fixed).

All numbers have a fixed 7 decimal places (18 digits total), and the maximum permitted value is +- 99999999999,
or just under 100 billion. NaN is supported.

The library is safe for concurrent use. Fixed values are immutable. It has built-in support for binary and json marshalling.

It is ideally suited for high performance trading financial systems. All common math operations are completed with 0 allocs.

**Design Goals**

Primarily developed to improve performance in [go-trader](https://github.com/robaho/go-trader).
Using Fixed rather than decimal.Decimal improves the performance by over 20%, and a lot less GC activity as well.
You can review these changes under the 'fixed' branch.

If you review the go-trader code, you will quickly see that I use dot imports for the fixed and common packages. Since this
is a "business/user" app and not systems code, this provides 2 major benefits: less verbose code, and I can easily change the
implementation of Fixed without changing lots of LOC - just the import statement, and some of the wrapper methods in common.

The fixed.Fixed API uses NaN for reporting errors in the common case, since often code is chained like:
```
   result := someFixed.Mul(NewS("123.50"))
```
and this would be a huge pain with error handling. Since all operations involving a NaN result in a NaN,
 any errors quickly surface anyway.

**Performance**

<pre>
using Go 1.25.3
cpu: Intel(R) Core(TM) i9-10850K CPU @ 3.60GHz
BenchmarkAddFixed-20            1000000000               0.7317 ns/op          0 B/op          0 allocs/op
BenchmarkAddDecimal-20          23494767                50.49 ns/op           80 B/op          2 allocs/op
BenchmarkAddBigInt-20           146968544                8.481 ns/op           0 B/op          0 allocs/op
BenchmarkAddBigFloat-20         19393024                59.52 ns/op           48 B/op          1 allocs/op
BenchmarkMulFixed-20            233557537                4.824 ns/op           0 B/op          0 allocs/op
BenchmarkMulDecimal-20          22572376                51.77 ns/op           80 B/op          2 allocs/op
BenchmarkMulBigInt-20           100000000               10.03 ns/op            0 B/op          0 allocs/op
BenchmarkMulBigFloat-20         56131700                21.59 ns/op            0 B/op          0 allocs/op
BenchmarkDivFixed-20            57347422                20.90 ns/op            0 B/op          0 allocs/op
BenchmarkDivDecimal-20           2935522               408.3 ns/op           384 B/op         12 allocs/op
BenchmarkDivBigInt-20           42439080                27.72 ns/op            8 B/op          1 allocs/op
BenchmarkDivBigFloat-20         14065242                84.40 ns/op            8 B/op          1 allocs/op
BenchmarkCmpFixed-20            1000000000               0.2113 ns/op          0 B/op          0 allocs/op
BenchmarkCmpDecimal-20          203959459                5.651 ns/op           0 B/op          0 allocs/op
BenchmarkCmpBigInt-20           311228709                3.796 ns/op           0 B/op          0 allocs/op
BenchmarkCmpBigFloat-20         319646821                3.756 ns/op           0 B/op          0 allocs/op
BenchmarkStringFixed-20         29299454                40.59 ns/op           24 B/op          1 allocs/op
BenchmarkStringNFixed-20        29791072                39.66 ns/op           24 B/op          1 allocs/op
BenchmarkStringDecimal-20        6768283               178.1 ns/op            56 B/op          4 allocs/op
BenchmarkStringBigInt-20        12304304                93.86 ns/op           16 B/op          1 allocs/op
BenchmarkStringBigFloat-20       3670106               325.5 ns/op           152 B/op          6 allocs/op
BenchmarkWriteTo-20             46491169                27.41 ns/op           23 B/op          0 allocs/op
</pre>

The "decimal" above is the common [shopspring decimal](https://github.com/shopspring/decimal) library

**Compatibility with SQL drivers**

By default `Fixed` implements `decomposer.Decimal` interface for database
drivers that support it. To use `sql.Scanner` and `driver.Valuer`
implementation flag `sql_scanner` must be specified on build.
