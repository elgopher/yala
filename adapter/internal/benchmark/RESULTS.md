# Test done in Jan 22 12:44AM

Please do not compare the results between different adapters, because they are testing different things. Results here
are useful for performance regression testing. 

* pkg: `github.com/jacekolszak/yala/adapter`
* cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz

```
BenchmarkGlog/global_logger_info-8         	 1208814	      1003 ns/op	     264 B/op	       4 allocs/op
BenchmarkGlog/local_logger_info-8          	 1213101	       992.8 ns/op	     264 B/op	       4 allocs/op
BenchmarkLog15/global_logger_info-8         	 1370102	       872.8 ns/op	     440 B/op	       3 allocs/op
BenchmarkLog15/local_logger_info-8          	 1210383	       987.1 ns/op	     440 B/op	       3 allocs/op
BenchmarkLogrus/global_logger_info-8         	  949262	      1232 ns/op	     404 B/op	      14 allocs/op
BenchmarkLogrus/local_logger_info-8          	 1000000	      1213 ns/op	     404 B/op	      14 allocs/op
BenchmarkPrinter/global_logger_info-8         	11838319	       118.0 ns/op	      40 B/op	       3 allocs/op
BenchmarkPrinter/local_logger_info-8          	14313648	        91.28 ns/op	      40 B/op	       3 allocs/op
BenchmarkZap/global_logger_info-8         	 7031860	       235.8 ns/op	     114 B/op	       2 allocs/op
BenchmarkZap/local_logger_info-8          	 5700342	       175.6 ns/op	     114 B/op	       2 allocs/op
BenchmarkZerolog/global_logger_info-8         	11077041	        99.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkZerolog/local_logger_info-8          	16146170	        73.07 ns/op	       0 B/op	       0 allocs/op
```