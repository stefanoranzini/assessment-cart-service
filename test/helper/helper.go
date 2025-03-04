package helper

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stefanoranzini/assessment/cart-service/internal/dao/db"
	"github.com/stefanoranzini/assessment/cart-service/tools/migration/migrate"
)

func PerpareTemporaryTestDb(t *testing.T) *sql.DB {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatal(err)
	}

	testDb, err := db.ConnectSQLite3(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = migrate.MigrateDB(testDb)
	if err != nil {
		t.Fatal(err)
	}

	return testDb
}
