package main

import (
	"log"

	"todo/cmd/api"
	"todo/config"
	"todo/storedb"
)

func main() {
	dbFilePath := storedb.GetDBFilePath()
	sqlStorage := storedb.NewSQLStorage(dbFilePath)

	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	store := storedb.NewStore(db)

	server := api.NewAPIServer(config.Envs.Port, store)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
