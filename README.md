
# JSON Database Evaluation

Contenders

1. ArangoDB 3.4.5 MMFile
1. ArangoDB 3.4.5 RocksDB
1. MongDB 4.1.11 WiredTiger
1. Couchbase Community 4.6.0
1. Postgres 11.3
1. MySQL 8.0 (yes, MySQL support JSON now)

## Performance Benchmark

### Run

Use Golang version >= 1.12

```
make run-bench
```

### Data structure

definition
```
{
    id: someguid,
    type: python|java|go|cpp,
    reference: someguid,
    version: 0.1.0,
    created_at: datetime,
    completed: count completed,
    failed: count failed,
    running: count running,
    elements: [
        {
            id:
            name:
        }
    ],
    content: string
}
```

instance
```
{
    id: someguid,
    type: python|java|golang|cpp,
    definition: someguid,
    reference: someguid,
    version: 0.1.0,
    created_at: datetime,
    status: {
        completed: false,
        failed: false,
        running: false,  
    },
    activities: [
        {
            id:
            name:
            created_at:
        }
    ]
}
```

### Query

See [main_test.go](main_test.go)

1. Insert single definitionn
1. Insert 100 definitions
1. Insert 100 definitions with each 200 instance
1. Get single definition
1. List definition
1. Get single definition with instance counter. See [QUERY.md](QUERY.md).

## Result

```
BenchmarkInsert1/mongodb-4         	    2000	    783665 
BenchmarkInsert1/arango-rock-4     	    1000	   1144089 
BenchmarkInsert1/arango-mm-4       	    2000	   1243355 
BenchmarkInsert1/postgres-4        	     500	   2707481 
BenchmarkInsert1/postgres11-4      	    1000	   1691399 
BenchmarkInsert1/mysql-4           	     500	   3921489 
BenchmarkInsert100/mongodb-4       	      20	  93714431 
BenchmarkInsert100/arango-rock-4   	      10	 128789873 
BenchmarkInsert100/arango-mm-4     	      20	  92419542 
BenchmarkInsert100/postgres-4      	      10	 195699543 
BenchmarkInsert100/postgres11-4    	      10	 185673680 
BenchmarkInsert100/mysql-4         	       3	 434513146 
BenchmarkGet/mongodb-100-4         	    1000	   1781482 
BenchmarkGet/mongodb-500-4         	    1000	   1567149 
BenchmarkGet/mongodb-1000-4        	    1000	   1839343 
BenchmarkGet/arango-rock-100-4     	    1000	   1453838 
BenchmarkGet/arango-rock-500-4     	    1000	   1121820 
BenchmarkGet/arango-rock-1000-4    	    1000	   1228395 
BenchmarkGet/arango-mm-100-4       	    2000	   1180700 
BenchmarkGet/arango-mm-500-4       	    1000	   1150089 
BenchmarkGet/arango-mm-1000-4      	    2000	   1070399 
BenchmarkGet/postgres-100-4        	     500	   2752511 
BenchmarkGet/postgres-500-4        	     500	   3075114 
BenchmarkGet/postgres-1000-4       	     500	   3194285 
BenchmarkGet/postgres11-100-4      	     500	   2997882 
BenchmarkGet/postgres11-500-4      	     500	   3667836 
BenchmarkGet/postgres11-1000-4     	     500	   2843938 
BenchmarkGet/mysql-100-4           	    1000	   1442468 
BenchmarkGet/mysql-500-4           	    1000	   1853118 
BenchmarkGet/mysql-1000-4          	    1000	   2320402 
BenchmarkList/mongodb-100-4        	     300	   4421753 
BenchmarkList/mongodb-500-4        	     200	   7593096 
BenchmarkList/mongodb-1000-4       	     200	   8438849 
BenchmarkList/arango-rock-100-4    	     300	   4304238 
BenchmarkList/arango-rock-500-4    	     200	   7391830 
BenchmarkList/arango-rock-1000-4   	     200	   8057644 
BenchmarkList/arango-mm-100-4      	     300	   4346995 
BenchmarkList/arango-mm-500-4      	     200	   9005103 
BenchmarkList/arango-mm-1000-4     	     200	   6994929 
BenchmarkList/postgres-100-4       	     300	   6032193 
BenchmarkList/postgres-500-4       	     100	  10141313 
BenchmarkList/postgres-1000-4      	     100	  12285722 
BenchmarkList/postgres11-100-4     	     300	   5916487 
BenchmarkList/postgres11-500-4     	     100	  10845324 
BenchmarkList/postgres11-1000-4    	     100	  11379720 
BenchmarkList/mysql-100-4          	     300	   4331380 
BenchmarkList/mysql-500-4          	     200	   8348725 
BenchmarkList/mysql-1000-4         	     200	  10488967 
BenchmarkInsertDI/mongodb-4        	      20	  85924235 
BenchmarkInsertDI/arango-rock-4    	      10	 129584884 
BenchmarkInsertDI/arango-mm-4      	      10	 112332012 
BenchmarkInsertDI/postgres-4       	      10	 222448863 
BenchmarkInsertDI/postgres11-4     	       5	 200868658 
BenchmarkInsertDI/mysql-4          	       2	 594101566 
BenchmarkGetStats/mongodb-100-4    	     500	   2540948 
BenchmarkGetStats/mongodb-500-4    	     500	   2693114 
BenchmarkGetStats/arango-rock-100-4         	    1000	   2403111 
BenchmarkGetStats/arango-rock-500-4         	     300	   3622294 
BenchmarkGetStats/arango-mm-100-4           	    1000	   1783055 
BenchmarkGetStats/arango-mm-500-4           	    1000	   3246957 
BenchmarkGetStats/postgres-100-4            	     200	   5690528 
BenchmarkGetStats/postgres-500-4            	     200	   5962688 
BenchmarkGetStats/postgres11-100-4          	     300	   5055073 
BenchmarkGetStats/postgres11-500-4          	     300	   5105726 
BenchmarkGetStats/mysql-100-4               	     500	   2347707 
BenchmarkGetStats/mysql-500-4               	     300	   4581015 
```