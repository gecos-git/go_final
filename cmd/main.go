package main

import (
	"fmt"
	"log"

	"todo/config"
	"todo/internal/api"
	"todo/internal/db"
	"todo/internal/store"
)

func main() {
	dbFilePath := db.GetDBFilePath()
	sqlStorage := db.NewSQLStorage(dbFilePath)

	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), store)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
