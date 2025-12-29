package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	db "github.com/nilesh0729/PixelScribe/internal/db/sqlc"
	"github.com/nilesh0729/PixelScribe/internal/api"
	"github.com/nilesh0729/PixelScribe/internal/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

