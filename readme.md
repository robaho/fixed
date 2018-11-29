**Summary**

A fixed place numeric library designed for performance.

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
BenchmarkAddFixed-8         	2000000000	         0.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddDecimal-8       	 3000000	       464 ns/op	     400 B/op	      10 allocs/op
BenchmarkAddBigInt-8        	100000000	        19.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddBigFloat-8      	20000000	       110 ns/op	      48 B/op	       1 allocs/op
BenchmarkMulFixed-8         	200000000	         6.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulDecimal-8       	20000000	        98.3 ns/op	      80 B/op	       2 allocs/op
BenchmarkMulBigInt-8        	100000000	        22.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulBigFloat-8      	30000000	        50.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDivFixed-8         	100000000	        20.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDivDecimal-8       	 1000000	      1189 ns/op	     928 B/op	      22 allocs/op
BenchmarkDivBigInt-8        	20000000	        72.8 ns/op	      48 B/op	       1 allocs/op
BenchmarkDivBigFloat-8      	10000000	       156 ns/op	      64 B/op	       2 allocs/op
BenchmarkCmpFixed-8         	2000000000	         0.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpDecimal-8       	100000000	        10.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpBigInt-8        	200000000	         8.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpBigFloat-8      	200000000	         8.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringFixed-8      	20000000	        97.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkStringDecimal-8    	 5000000	       335 ns/op	     144 B/op	       5 allocs/op
BenchmarkStringBigInt-8     	10000000	       215 ns/op	      80 B/op	       3 allocs/op
BenchmarkStringBigFloat-8   	 3000000	       586 ns/op	     272 B/op	       8 allocs/op
</pre>

The "decimal" above is the common [shopspring decimal](https://github.com/shopspring/decimal) library