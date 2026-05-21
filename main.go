package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Request struct {
	ID        int
	IP        string
	Method    string
	Path      string
	Timestamp string
}

func main() {
	connStr := "host=localhost port=5433 user=postgres password=12345678 dbname=mirage sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mirage is live.")

	store := NewPostgresStore(db)
	server := NewServer(store)
	server.Start("8080")
}
