**Summary**

A fixed place numeric library designed for performance.

The c++ version is available [here](https://github.com/robaho/cpp_fixed).

All numbers have a fixed 7 decimal places, and the maximum permitted value is +- 99999999999,
or just under 100 billion.

The library is safe for concurrent use. It has built-in support for binary and json marshalling.

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
using Go 1.21.5
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkAddFixed-8             1000000000               0.9627 ns/op          0 B/op          0 allocs/op
BenchmarkAddDecimal-8           17871763                66.52 ns/op           80 B/op          2 allocs/op
BenchmarkAddBigInt-8            125826048                9.562 ns/op           0 B/op          0 allocs/op
BenchmarkAddBigFloat-8          18763552                63.51 ns/op           48 B/op          1 allocs/op
BenchmarkMulFixed-8             335886367                3.538 ns/op           0 B/op          0 allocs/op
BenchmarkMulDecimal-8           18164803                66.12 ns/op           80 B/op          2 allocs/op
BenchmarkMulBigInt-8            100000000               10.41 ns/op            0 B/op          0 allocs/op
BenchmarkMulBigFloat-8          50151100                23.93 ns/op            0 B/op          0 allocs/op
BenchmarkDivFixed-8             328157694                3.722 ns/op           0 B/op          0 allocs/op
BenchmarkDivDecimal-8            2558497               461.7 ns/op           384 B/op         12 allocs/op
BenchmarkDivBigInt-8            33726384                34.68 ns/op            8 B/op          1 allocs/op
BenchmarkDivBigFloat-8          10757650               110.1 ns/op            24 B/op          2 allocs/op
BenchmarkCmpFixed-8             1000000000               0.2519 ns/op          0 B/op          0 allocs/op
BenchmarkCmpDecimal-8           171236422                6.926 ns/op           0 B/op          0 allocs/op
BenchmarkCmpBigInt-8            250970304                4.791 ns/op           0 B/op          0 allocs/op
BenchmarkCmpBigFloat-8          271898336                4.428 ns/op           0 B/op          0 allocs/op
BenchmarkStringFixed-8          23637406                50.30 ns/op           24 B/op          1 allocs/op
BenchmarkStringNFixed-8         23457960                51.85 ns/op           24 B/op          1 allocs/op
BenchmarkStringDecimal-8         5763308               210.2 ns/op            56 B/op          4 allocs/op
BenchmarkStringBigInt-8         11742596               114.0 ns/op            16 B/op          1 allocs/op
BenchmarkStringBigFloat-8        3003280               395.3 ns/op           176 B/op          7 allocs/op
BenchmarkWriteTo-8              38573978                43.13 ns/op           27 B/op          0 allocs/op
</pre>

The "decimal" above is the common [shopspring decimal](https://github.com/shopspring/decimal) library

**Compatibility with SQL drivers**

By default `Fixed` implements `decomposer.Decimal` interface for database
drivers that support it. To use `sql.Scanner` and `driver.Valuer`
implementation flag `sql_scanner` must be specified on build.
