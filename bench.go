package prepared_stmt_benchmark

import "database/sql"

func Connect() (*sql.DB, error) {
	return sql.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
}

func CreateTestTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS test (
				id INT NOT NULL AUTO_INCREMENT,
				data VARCHAR(255) NOT NULL,
				PRIMARY KEY (id)
			)`)
	return err
}

func DropTestTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE test")
	return err
}
