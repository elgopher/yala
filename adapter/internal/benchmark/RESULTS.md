# Test done in Jan 22 12:44AM

Please do not compare the results between different adapters, because they are testing different things. Results here
are useful for performance regression testing. 

* cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz

```
BenchmarkConsole/global_logger_info-8         	 8403124	       140.6 ns/op	      40 B/op	       3 allocs/op
BenchmarkConsole/local_logger_info-8          	 9990122	       123.8 ns/op	      40 B/op	       3 allocs/op
BenchmarkConsole/time/local_logger_with_field-8         	 1761192	       691.3 ns/op	     200 B/op	       8 allocs/op
BenchmarkConsole/string/local_logger_with_field-8       	 4456893	       297.5 ns/op	      96 B/op	       6 allocs/op
BenchmarkConsole/int/local_logger_with_field-8          	 4396735	       279.8 ns/op	      88 B/op	       5 allocs/op
BenchmarkConsole/int64/local_logger_with_field-8        	 4330046	       277.9 ns/op	      96 B/op	       6 allocs/op
BenchmarkConsole/float64/local_logger_with_field-8      	 3428498	       357.2 ns/op	      96 B/op	       6 allocs/op
BenchmarkConsole/float32/local_logger_with_field-8      	 3494066	       345.2 ns/op	      96 B/op	       6 allocs/op
BenchmarkGlog/global_logger_info-8         	 1000000	      1131 ns/op	     232 B/op	       3 allocs/op
BenchmarkGlog/local_logger_info-8          	 1000000	      1012 ns/op	     232 B/op	       3 allocs/op
BenchmarkGlog/int64/local_logger_with_field-8         	  932733	      1293 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/float64/local_logger_with_field-8       	  805698	      1430 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/float32/local_logger_with_field-8       	  882258	      1406 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/time/local_logger_with_field-8          	  693415	      1776 ns/op	     448 B/op	      10 allocs/op
BenchmarkGlog/string/local_logger_with_field-8        	  960004	      1265 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/int/local_logger_with_field-8           	  956354	      1266 ns/op	     288 B/op	       6 allocs/op
BenchmarkLog15/global_logger_info-8         	  814569	      1626 ns/op	     680 B/op	      12 allocs/op
BenchmarkLog15/local_logger_info-8          	  722480	      1722 ns/op	     680 B/op	      12 allocs/op
BenchmarkLog15/float32/local_logger_with_field-8         	  536744	      2173 ns/op	    1109 B/op	      23 allocs/op
BenchmarkLog15/time/local_logger_with_field-8            	  587264	      2520 ns/op	    1264 B/op	      23 allocs/op
BenchmarkLog15/string/local_logger_with_field-8          	  612498	      2306 ns/op	    1080 B/op	      21 allocs/op
BenchmarkLog15/int/local_logger_with_field-8             	  536482	      2097 ns/op	    1080 B/op	      21 allocs/op
BenchmarkLog15/int64/local_logger_with_field-8           	  579112	      2419 ns/op	    1082 B/op	      22 allocs/op
BenchmarkLog15/float64/local_logger_with_field-8         	  556372	      2568 ns/op	    1109 B/op	      23 allocs/op
BenchmarkLog/global_logger_info-8         	 4135249	       285.2 ns/op	      56 B/op	       4 allocs/op
BenchmarkLog/local_logger_info-8          	 4483218	       266.0 ns/op	      56 B/op	       4 allocs/op
BenchmarkLog/time/local_logger_with_field-8         	 1414213	       829.7 ns/op	     248 B/op	       9 allocs/op
BenchmarkLog/string/local_logger_with_field-8       	 2831449	       424.3 ns/op	     112 B/op	       7 allocs/op
BenchmarkLog/int/local_logger_with_field-8          	 2871441	       422.4 ns/op	     104 B/op	       6 allocs/op
BenchmarkLog/int64/local_logger_with_field-8        	 2852680	       425.4 ns/op	     112 B/op	       7 allocs/op
BenchmarkLog/float64/local_logger_with_field-8      	 2423032	       491.6 ns/op	     112 B/op	       7 allocs/op
BenchmarkLog/float32/local_logger_with_field-8      	 2464716	       484.3 ns/op	     112 B/op	       7 allocs/op
BenchmarkLogrus/global_logger_info-8         	  978819	      1195 ns/op	     404 B/op	      14 allocs/op
BenchmarkLogrus/local_logger_info-8          	 1000000	      1143 ns/op	     404 B/op	      14 allocs/op
BenchmarkLogrus/time/local_logger_with_field-8         	  472927	      2988 ns/op	    1316 B/op	      24 allocs/op
BenchmarkLogrus/string/local_logger_with_field-8       	  482336	      2285 ns/op	    1204 B/op	      20 allocs/op
BenchmarkLogrus/int/local_logger_with_field-8          	  626422	      2337 ns/op	    1204 B/op	      20 allocs/op
BenchmarkLogrus/int64/local_logger_with_field-8        	  588906	      2429 ns/op	    1208 B/op	      21 allocs/op
BenchmarkLogrus/float64/local_logger_with_field-8      	  548394	      2435 ns/op	    1208 B/op	      21 allocs/op
BenchmarkLogrus/float32/local_logger_with_field-8      	  643100	      2115 ns/op	    1208 B/op	      21 allocs/op
BenchmarkZap/global_logger_info-8         	 5978613	       214.4 ns/op	     114 B/op	       2 allocs/op
BenchmarkZap/local_logger_info-8          	 7606512	       182.7 ns/op	     114 B/op	       2 allocs/op
BenchmarkZap/string/local_logger_with_field-8         	 4325899	       325.1 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/int/local_logger_with_field-8            	 4670941	       343.6 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/int64/local_logger_with_field-8          	 3283704	       316.4 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/float64/local_logger_with_field-8        	 2924810	       369.9 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/float32/local_logger_with_field-8        	 4364494	       334.7 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/time/local_logger_with_field-8           	 3555284	       393.8 ns/op	     234 B/op	       5 allocs/op
BenchmarkZerolog/global_logger_info-8         	10537152	        97.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkZerolog/local_logger_info-8          	15246154	        78.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkZerolog/string/local_logger_with_field-8         	 8554328	       139.8 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/int/local_logger_with_field-8            	 8463698	       139.3 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/int64/local_logger_with_field-8          	 8701618	       137.0 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/float64/local_logger_with_field-8        	 5781765	       207.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/float32/local_logger_with_field-8        	 5923426	       202.5 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/time/local_logger_with_field-8           	 4157388	       287.6 ns/op	      32 B/op	       1 allocs/op
BenchmarkInfo-8          	38938502	        30.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogger_Info-8   	100000000	        10.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkWith-8          	30173617	        48.44 ns/op	      32 B/op	       1 allocs/op
BenchmarkWithError-8     	91293628	        12.95 ns/op	       0 B/op	       0 allocs/op
```