# DB Performance Benchmark

This is a quick benchmark of running common queries against a database using different methods. We will be testing out
how GORM compares with the standard library and then also comparing the performance of using a prepared statement over 
regular queries.

## Setup

We have a simple database with a single table called `test` with the following schema:

```sql
CREATE TABLE IF NOT EXISTS test (
    id INT NOT NULL AUTO_INCREMENT,
    data VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);
```

* For the first test, we will be executing a simple `INSERT` query to insert a single row into the table. 
* For the second test, we will be fetching a single row from the table using the `id` column.
* For the third test and the final, we will be fetching all rows from the table.

## Results

For this benchmark, we have considered the performance of the `database/sql` library paired with the `github.com/go-sql-driver/mysql` driver
as the baseline. We have then measured the change in performance when using GORM or a PreparedStatement for the same query. Here are the results:

```bash
➜  prepared-stmt-benchmark git:(main) ✗ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: prepared-stmt-benchmark
BenchmarkNonPreparedStmtExec-10                      782           1329609 ns/op             320 B/op         13 allocs/op
BenchmarkNonPreparedStmtExecGorm-10                  552           2042676 ns/op            5511 B/op        107 allocs/op
BenchmarkPreparedStmtExec-10                        1268            971607 ns/op             248 B/op         11 allocs/op
BenchmarkNonPreparedStmtQueryRow-10                 1731            734764 ns/op             859 B/op         27 allocs/op
BenchmarkNonPreparedStmtQueryRowGorm-10             1545            767889 ns/op           11087 B/op        203 allocs/op
BenchmarkPreparedStmtQueryRow-10                    3184            380564 ns/op             739 B/op         25 allocs/op
BenchmarkNonPreparedStmtQuery-10                     808           1577949 ns/op           56122 B/op       1959 allocs/op
BenchmarkNonPreparedStmtQueryGorm-10                 465           3481287 ns/op          872735 B/op      14686 allocs/op
BenchmarkPreparedStmtQuery-10                       1298           1085053 ns/op           55955 B/op       1953 allocs/op
PASS
ok      prepared-stmt-benchmark 44.767s
````

### Exec Query

* **Baseline** - This query is executed 782 times in the duration of the benchmark consuming 320 bytes of memory.
* **GORM** - This query is executed 552 times in the duration of the benchmark consuming 5511 bytes of memory. This is
about 29% slower than the baseline while consuming 17 times more memory!
* **Prepared Statement** - This query is executed 1268 times in the duration of the benchmark consuming 248 bytes of memory.
This is about 62% faster than the baseline while consuming 22% less memory!

### QueryRow

* **Baseline** - This query is executed 1731 times in the duration of the benchmark consuming 859 bytes of memory.
* **GORM** - This query is executed 1545 times in the duration of the benchmark consuming 11087 bytes of memory. This is
about 11% slower than the baseline while consuming 12 times more memory!
* **Prepared Statement** - This query is executed 3184 times in the duration of the benchmark consuming 739 bytes of memory.
This is about 84% faster than the baseline while consuming 14% less memory!

### Query

* **Baseline** - This query is executed 808 times in the duration of the benchmark consuming 56.122 kB of memory.
* **GORM** - This query is executed 465 times in the duration of the benchmark consuming 872.735 kB of memory. This is
about 42% slower than the baseline while consuming 15 times more memory!
* **Prepared Statement** - This query is executed 1298 times in the duration of the benchmark consuming 55.955 kB of memory.
* This is about 61% faster than the baseline while consuming 0.9% less memory!


## Conclusion

As we can see from the results, using a prepared statement is significantly faster than using a regular query. This is
because the prepared statement is compiled once and then executed multiple times. This is especially useful when we are
executing the same query multiple times with different parameters.

## TLDR 

1. DO NOT USE GORM FOR PERFORMANCE CRITICAL APPLICATIONS!
2. USE PREPARED STATEMENTS WHENEVER POSSIBLE!