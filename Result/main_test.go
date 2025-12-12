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

func TestMain(m *testing.M) {
	conn, err := sql.Open(DbDriver, DbSource)
	if err != nil{
		log.Fatal("cannot connect to DB: ", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
