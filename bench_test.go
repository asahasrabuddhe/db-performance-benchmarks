package prepared_stmt_benchmark

import (
	"database/sql"
	"errors"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkNonPreparedStmtExec(b *testing.B) {
	db, err := Connect()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tt := Test{Data: uuid.New().String()}
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", tt.Data)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = DropTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNonPreparedStmtExecGorm(b *testing.B) {
	db, err := ConnectGORM()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTableGORM(db, &Test{})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tt := Test{Data: uuid.New().String()}
		err = db.Save(&tt).Error
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = DropTestTableGORM(db, &Test{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkPreparedStmtExec(b *testing.B) {
	db, err := Connect()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO test (data) VALUES (?)")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tt := Test{Data: uuid.New().String()}
		_, err = stmt.Exec(tt.Data)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = stmt.Close()
	if err != nil {
		b.Fatal(err)
	}

	err = DropTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNonPreparedStmtQueryRow(b *testing.B) {
	db, err := Connect()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		tt := Test{Data: uuid.New().String()}
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", tt.Data)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var tt Test
		err = db.QueryRow("SELECT * FROM test WHERE id = ?", rand.Int31n(1000)).Scan(&tt.ID, &tt.Data)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = DropTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNonPreparedStmtQueryRowGorm(b *testing.B) {
	db, err := ConnectGORM()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTableGORM(db, &Test{})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		tt := Test{Data: uuid.New().String()}
		err = db.Save(&tt).Error
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var tt Test
		err = db.Where("id = ?", rand.Int31n(1000)).First(&tt).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				continue
			}
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = DropTestTableGORM(db, &Test{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkPreparedStmtQueryRow(b *testing.B) {
	db, err := Connect()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		tt := Test{Data: uuid.New().String()}
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", tt.Data)
		if err != nil {
			b.Fatal(err)
		}
	}

	stmt, err := db.Prepare("SELECT * FROM test WHERE id = ?")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var tt Test
		err = stmt.QueryRow(rand.Int31n(1000)).Scan(&tt.ID, &tt.Data)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = stmt.Close()
	if err != nil {
		b.Fatal(err)
	}

	err = DropTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNonPreparedStmtQuery(b *testing.B) {
	db, err := Connect()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		tt := Test{Data: uuid.New().String()}
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", tt.Data)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rowCount := rand.Int31n(1000)
		var rows *sql.Rows
		rows, err = db.Query("SELECT * FROM test WHERE id > ? ORDER BY RAND()", rowCount)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
		ttt := make([]Test, 1000-rowCount)
		var index int
		for rows.Next() {
			err = rows.Scan(&ttt[index].ID, &ttt[index].Data)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}
				b.Fatal(err)
			}
			index++
		}
	}

	b.StopTimer()

	err = DropTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNonPreparedStmtQueryGorm(b *testing.B) {
	db, err := ConnectGORM()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTableGORM(db, &Test{})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		tt := Test{Data: uuid.New().String()}
		err = db.Save(&tt).Error
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rowCount := rand.Int31n(1000)
		tt := make([]Test, 1000-rowCount)
		err = db.Where("id > ?", rowCount).Find(&tt).Error
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
	}

	b.StopTimer()

	err = DropTestTableGORM(db, &Test{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkPreparedStmtQuery(b *testing.B) {
	db, err := Connect()
	if err != nil {
		b.Fatal(err)
	}

	err = CreateTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		tt := Test{Data: uuid.New().String()}
		_, err = db.Exec("INSERT INTO test (data) VALUES (?)", tt.Data)
		if err != nil {
			b.Fatal(err)
		}
	}

	stmt, err := db.Prepare("SELECT * FROM test WHERE id > ? ORDER BY RAND()")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rowCount := rand.Int31n(1000)
		var rows *sql.Rows
		rows, err = stmt.Query(rowCount)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			b.Fatal(err)
		}
		ttt := make([]Test, 1000-rowCount)
		var index int
		for rows.Next() {
			err = rows.Scan(&ttt[index].ID, &ttt[index].Data)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}
				b.Fatal(err)
			}
			index++
		}
	}

	b.StopTimer()

	err = stmt.Close()
	if err != nil {
		b.Fatal(err)
	}

	err = DropTestTable(db)
	if err != nil {
		b.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		b.Fatal(err)
	}
}
