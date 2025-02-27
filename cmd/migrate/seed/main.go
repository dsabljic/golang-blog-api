package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewPostgresStorage(conn)

	db.Seed(store, conn)
}
