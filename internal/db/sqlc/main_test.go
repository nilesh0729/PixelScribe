package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_"github.com/lib/pq"
)

const (
	DbDriver = "postgres"
	DbSource = "postgres://root:secret@localhost:5434/Pros?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	conn, err := sql.Open(DbDriver, DbSource)
	if err != nil{
		log.Fatal("cannot connect to DB: ", err)
	}
	testDB = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}
