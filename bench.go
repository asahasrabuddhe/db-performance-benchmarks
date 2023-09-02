package prepared_stmt_benchmark

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/projectleo/gorm"
	_ "github.com/projectleo/gorm/dialects/mysql"
)

type Test struct {
	ID   int
	Data string
}

func Connect() (*sql.DB, error) {
	return sql.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
}

func ConnectGORM() (*gorm.DB, error) {
	return gorm.Open("mysql", "benchuser:benchp@ss@tcp(localhost:3306)/benchdb")
}

func CreateTestTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS test (
				id INT NOT NULL AUTO_INCREMENT,
				data VARCHAR(255) NOT NULL,
				PRIMARY KEY (id)
			)`)
	return err
}

func CreateTestTableGORM(db *gorm.DB, model any) error {
	return db.AutoMigrate(model).Error
}

func DropTestTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE test")
	return err
}

func DropTestTableGORM(db *gorm.DB, model any) error {
	return db.DropTable(model).Error
}
