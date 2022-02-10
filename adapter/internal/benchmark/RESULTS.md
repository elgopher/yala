# Test done in Jan 22 12:44AM

Please do not compare the results between different adapters, because they are testing different things. Results here
are useful for performance regression testing. 

* cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz

```
BenchmarkConsole/global_logger_info-8         	11107899	       106.1 ns/op	      24 B/op	       2 allocs/op
BenchmarkConsole/normal_logger_info-8         	12160440	        98.95 ns/op	      24 B/op	       2 allocs/op
BenchmarkConsole/int/normal_logger_with_field-8         	 4782313	       255.9 ns/op	      72 B/op	       4 allocs/op
BenchmarkConsole/int64/normal_logger_with_field-8       	 4588833	       258.3 ns/op	      80 B/op	       5 allocs/op
BenchmarkConsole/float64/normal_logger_with_field-8     	 3689180	       328.6 ns/op	      80 B/op	       5 allocs/op
BenchmarkConsole/float32/normal_logger_with_field-8     	 3720772	       320.1 ns/op	      80 B/op	       5 allocs/op
BenchmarkConsole/time/normal_logger_with_field-8        	 1815115	       667.6 ns/op	     184 B/op	       7 allocs/op
BenchmarkConsole/string/normal_logger_with_field-8      	 4461062	       256.1 ns/op	      80 B/op	       5 allocs/op
BenchmarkGlog/global_logger_info-8         	 1000000	      1069 ns/op	     232 B/op	       3 allocs/op
BenchmarkGlog/normal_logger_info-8         	 1000000	      1040 ns/op	     232 B/op	       3 allocs/op
BenchmarkGlog/time/normal_logger_with_field-8         	  695626	      1794 ns/op	     448 B/op	      10 allocs/op
BenchmarkGlog/string/normal_logger_with_field-8       	  940095	      1276 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/int/normal_logger_with_field-8          	  957658	      1258 ns/op	     288 B/op	       6 allocs/op
BenchmarkGlog/int64/normal_logger_with_field-8        	  906349	      1319 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/float64/normal_logger_with_field-8      	  863223	      1390 ns/op	     296 B/op	       7 allocs/op
BenchmarkGlog/float32/normal_logger_with_field-8      	  862776	      1390 ns/op	     296 B/op	       7 allocs/op
BenchmarkLog15/global_logger_info-8         	  805639	      1679 ns/op	     680 B/op	      12 allocs/op
BenchmarkLog15/normal_logger_info-8         	  874143	      1699 ns/op	     680 B/op	      12 allocs/op
BenchmarkLog15/int64/normal_logger_with_field-8         	  541711	      2067 ns/op	    1082 B/op	      22 allocs/op
BenchmarkLog15/float64/normal_logger_with_field-8       	  592009	      2395 ns/op	    1109 B/op	      23 allocs/op
BenchmarkLog15/float32/normal_logger_with_field-8       	  579549	      2191 ns/op	    1109 B/op	      23 allocs/op
BenchmarkLog15/time/normal_logger_with_field-8          	  534758	      2655 ns/op	    1264 B/op	      23 allocs/op
BenchmarkLog15/string/normal_logger_with_field-8        	  540922	      2290 ns/op	    1080 B/op	      21 allocs/op
BenchmarkLog15/int/normal_logger_with_field-8           	  652357	      2183 ns/op	    1080 B/op	      21 allocs/op
BenchmarkLog/global_logger_info-8         	 6773872	       169.2 ns/op	       8 B/op	       1 allocs/op
BenchmarkLog/normal_logger_info-8         	 7405096	       161.7 ns/op	       8 B/op	       1 allocs/op
BenchmarkLog/int/normal_logger_with_field-8         	 3841082	       310.7 ns/op	      56 B/op	       3 allocs/op
BenchmarkLog/int64/normal_logger_with_field-8       	 3771975	       318.9 ns/op	      64 B/op	       4 allocs/op
BenchmarkLog/float64/normal_logger_with_field-8     	 3059952	       398.9 ns/op	      64 B/op	       4 allocs/op
BenchmarkLog/float32/normal_logger_with_field-8     	 3073831	       391.5 ns/op	      64 B/op	       4 allocs/op
BenchmarkLog/time/normal_logger_with_field-8        	 1667187	       721.1 ns/op	     168 B/op	       6 allocs/op
BenchmarkLog/string/normal_logger_with_field-8      	 3767346	       314.3 ns/op	      64 B/op	       4 allocs/op
BenchmarkLogrus/global_logger_info-8         	  913322	      1333 ns/op	     468 B/op	      16 allocs/op
BenchmarkLogrus/normal_logger_info-8                           	  879543	      1298 ns/op	     468 B/op	      16 allocs/op
BenchmarkLogrus/string/normal_logger_with_field-8              	  514377	      3583 ns/op	    1605 B/op	      24 allocs/op
BenchmarkLogrus/int/normal_logger_with_field-8                 	  309808	      3697 ns/op	    1605 B/op	      24 allocs/op
BenchmarkLogrus/int64/normal_logger_with_field-8               	  465054	      4021 ns/op	    1609 B/op	      25 allocs/op
BenchmarkLogrus/float64/normal_logger_with_field-8             	  255966	      3953 ns/op	    1609 B/op	      25 allocs/op
BenchmarkLogrus/float32/normal_logger_with_field-8             	  346702	      3924 ns/op	    1609 B/op	      25 allocs/op
BenchmarkLogrus/time/normal_logger_with_field-8                	  414996	      3643 ns/op	    1717 B/op	      28 allocs/op
BenchmarkLogrus/global_logger_info_with_three_fields-8         	  443599	      3990 ns/op	    1733 B/op	      24 allocs/op
BenchmarkLogrus/normal_logger_error_with_cause_and_two_fields-8         	  362486	      3647 ns/op	    1713 B/op	      25 allocs/op
BenchmarkZap/global_logger_info-8         	 5975440	       189.1 ns/op	     114 B/op	       2 allocs/op
BenchmarkZap/normal_logger_info-8         	 8121460	       196.7 ns/op	     114 B/op	       2 allocs/op
BenchmarkZap/float64/normal_logger_with_field-8         	 4922467	       306.0 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/float32/normal_logger_with_field-8         	 4165966	       339.0 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/time/normal_logger_with_field-8            	 3853836	       384.2 ns/op	     234 B/op	       5 allocs/op
BenchmarkZap/string/normal_logger_with_field-8          	 4434114	       382.0 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/int/normal_logger_with_field-8             	 3943224	       344.7 ns/op	     210 B/op	       4 allocs/op
BenchmarkZap/int64/normal_logger_with_field-8           	 4745278	       339.3 ns/op	     210 B/op	       4 allocs/op
BenchmarkZerolog/global_logger_info-8         	14133230	        82.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkZerolog/normal_logger_info-8         	14850610	        78.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkZerolog/int/normal_logger_with_field-8         	 8845273	       134.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/int64/normal_logger_with_field-8       	 8795118	       134.3 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/float64/normal_logger_with_field-8     	 5896245	       204.6 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/float32/normal_logger_with_field-8     	 6067651	       198.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/time/normal_logger_with_field-8        	 4251734	       280.0 ns/op	      32 B/op	       1 allocs/op
BenchmarkZerolog/string/normal_logger_with_field-8      	 8692183	       137.5 ns/op	      32 B/op	       1 allocs/op
BenchmarkInfo-8           	97691961	        12.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogger_Info-8    	93976857	        12.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkWith-8           	40131960	        38.78 ns/op	      32 B/op	       1 allocs/op
BenchmarkWithError-8      	450374347	         2.599 ns/op	       0 B/op	       0 allocs/op
```