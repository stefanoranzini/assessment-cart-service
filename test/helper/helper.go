package helper

import (
	"database/sql"
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao/db"
	"github.com/stefanoranzini/assessment/cart-service/tools/migration/migrate"
)

func PerpareTemporaryTestDb(t *testing.T) *sql.DB {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatal(err)
	}

	return PrepareTestDb(t, tempFile.Name())
}

func PrepareTestDb(t *testing.T, dataSourceName string) *sql.DB {
	db, err := db.ConnectSQLite3(dataSourceName)
	if err != nil {
		t.Fatal(err)
	}

	err = migrate.MigrateDB(db)
	if err != nil {
		t.Fatal(err)
	}

	// For assessment purpose some data are inserted via migrations (not a common production scenario), but for testing we want to control the data present on the tables
	TruncateProductTable(t, db)

	return db
}

func InsertProduct(t *testing.T, db *sql.DB, id int, price decimal.Decimal) {
	_, err := db.Exec("INSERT INTO product (id, price) VALUES (?, ?)", id, price)
	if err != nil {
		t.Fatal(err)
	}
}

func TruncateProductTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM product")
	if err != nil {
		t.Fatal(err)
	}
}
