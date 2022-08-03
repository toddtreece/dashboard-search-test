```
go test -bench=. -benchmem | benchstat /dev/stdin
name                          time/op
SQL/regex_100-32               728µs ± 0%
SQL/regex_1000-32             8.36ms ± 0%
SQL/regex_10000-32            85.2ms ± 0%
SQL/sql_search_100-32          325µs ± 0%
SQL/sql_search_1000-32        1.47ms ± 0%
SQL/sql_search_10000-32       13.4ms ± 0%
SQL/bluge_search_100-32       1.04ms ± 0%
SQL/bluge_search_1000-32      1.36ms ± 0%
SQL/bluge_search_10000-32     4.98ms ± 0%
Index/bluge_index_100-32      72.3ms ± 0%
Index/bluge_index_1000-32      2.28s ± 0%
Index/bluge_index_10000-32     11.7s ± 0%
Index/sql_db_create_100-32    4.95ms ± 0%
Index/sql_db_create_1000-32   53.1ms ± 0%
Index/sql_db_create_10000-32   543ms ± 0%

name                          alloc/op
SQL/regex_100-32               431kB ± 0%
SQL/regex_1000-32             4.31MB ± 0%
SQL/regex_10000-32            43.2MB ± 0%
SQL/sql_search_100-32          114kB ± 0%
SQL/sql_search_1000-32         491kB ± 0%
SQL/sql_search_10000-32       4.23MB ± 0%
SQL/bluge_search_100-32        195kB ± 0%
SQL/bluge_search_1000-32       247kB ± 0%
SQL/bluge_search_10000-32     1.40MB ± 0%
Index/bluge_index_100-32       122MB ± 0%
Index/bluge_index_1000-32     1.76GB ± 0%
Index/bluge_index_10000-32    9.39GB ± 0%
Index/sql_db_create_100-32     593kB ± 0%
Index/sql_db_create_1000-32   5.04MB ± 0%
Index/sql_db_create_10000-32  50.3MB ± 0%

name                          allocs/op
SQL/regex_100-32                 937 ± 0%
SQL/regex_1000-32              9.11k ± 0%
SQL/regex_10000-32             92.3k ± 0%
SQL/sql_search_100-32          1.74k ± 0%
SQL/sql_search_1000-32         8.06k ± 0%
SQL/sql_search_10000-32        71.4k ± 0%
SQL/bluge_search_100-32          403 ± 0%
SQL/bluge_search_1000-32       1.38k ± 0%
SQL/bluge_search_10000-32      18.2k ± 0%
Index/bluge_index_100-32        148k ± 0%
Index/bluge_index_1000-32      12.9M ± 0%
Index/bluge_index_10000-32     55.4M ± 0%
Index/sql_db_create_100-32     2.85k ± 0%
Index/sql_db_create_1000-32    26.3k ± 0%
Index/sql_db_create_10000-32    270k ± 0%
```