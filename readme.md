**Summary**

A fixed place numeric library designed for performance.

All numbers have a fixed 7 decimal places, and the maximum permitted value is +- 99999999999,
or just under 100 billion.

It is ideally suited for high performance trading financial systems. All common math operations are completed with 0 allocs.

**Performance**

```
BenchmarkAddFixed-8         	2000000000	         0.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddDecimal-8       	 3000000	       457 ns/op	     400 B/op	      10 allocs/op
BenchmarkAddBigInt-8        	100000000	        19.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddBigFloat-8      	20000000	       110 ns/op	      48 B/op	       1 allocs/op
BenchmarkMulFixed-8         	100000000	        12.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulDecimal-8       	20000000	        94.2 ns/op	      80 B/op	       2 allocs/op
BenchmarkMulBigInt-8        	100000000	        22.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulBigFloat-8      	30000000	        50.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDivFixed-8         	100000000	        19.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDivDecimal-8       	 1000000	      1152 ns/op	     928 B/op	      22 allocs/op
BenchmarkDivBigInt-8        	20000000	        68.4 ns/op	      48 B/op	       1 allocs/op
BenchmarkDivBigFloat-8      	10000000	       151 ns/op	      64 B/op	       2 allocs/op
BenchmarkCmpFixed-8         	2000000000	         0.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpDecimal-8       	100000000	        10.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpBigInt-8        	200000000	         8.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpBigFloat-8      	200000000	         7.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringFixed-8      	20000000	        99.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkStringDecimal-8    	 5000000	       326 ns/op	     144 B/op	       5 allocs/op
BenchmarkStringBigInt-8     	10000000	       209 ns/op	      80 B/op	       3 allocs/op
BenchmarkStringBigFloat-8   	 3000000	       571 ns/op	     272 B/op	       8 allocs/op
```

The "decimal" above is the common [shopspring decimal](https://github.com/shopspring/decimal) library