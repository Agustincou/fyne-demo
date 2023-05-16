package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

const sqlPath = "./testdata/sql.db"

var testSqliteRepo *SQLiteRepository

func TestMain(m *testing.M) {
	_ = os.Remove(sqlPath)
	db, err := sql.Open("sqlite", sqlPath)
	if err != nil {
		log.Println(err)
	}

	testSqliteRepo = NewSQLiteRespository(db)
	os.Exit(m.Run())
}
