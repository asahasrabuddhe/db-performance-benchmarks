package prepared_stmt_benchmark

import (
	"database/sql"
	"errors"
	"math/rand"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func BenchmarkNonPreparedStmtExec(b *testing.B) {
	db, err := sql.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (
				id INT NOT NULL AUTO_INCREMENT,
				data VARCHAR(255) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", uuid.New().String())
		if err != nil {
			b.Fatal(err)
		}
	}

	_, err = db.Exec("DROP TABLE test")
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkPreparedStmtExec(b *testing.B) {
	db, err := sql.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (
				id INT NOT NULL AUTO_INCREMENT,
				data VARCHAR(255) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		b.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO test (data) VALUES (?)")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = stmt.Exec(uuid.New().String())
		if err != nil {
			b.Fatal(err)
		}
	}

	err = stmt.Close()
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE test")
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNonPreparedStmtQuery(b *testing.B) {
	db, err := sql.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (
				id INT NOT NULL AUTO_INCREMENT,
				data VARCHAR(255) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", uuid.New().String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var data string
		err = db.QueryRow("SELECT data FROM test WHERE id = ?", rand.Int31n(1000)).Scan(&data)
		if err != nil {
			b.Fatal(err)
		}
	}

	_, err = db.Exec("DROP TABLE test")
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkPreparedStmtQuery(b *testing.B) {
	db, err := sql.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (
				id INT NOT NULL AUTO_INCREMENT,
				data VARCHAR(255) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", uuid.New().String())
		if err != nil {
			b.Fatal(err)
		}
	}

	stmt, err := db.Prepare("SELECT data FROM test WHERE id = ?")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var data string
		err = stmt.QueryRow(rand.Int31n(1000)).Scan(&data)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	err = stmt.Close()
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE test")
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}
