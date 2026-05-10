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

func run(store Store) {
	store.Insert("10.10.10.1", "GET", "/secret")
	store.Insert("10.10.10.2", "POST", "/admin/login")

	fmt.Println("Inserted 2 requests via interface.")

	requests := store.GetAll()

	fmt.Printf("\n%-5s %-15s %-8s %-20s %s\n", "ID", "IP", "Method", "Path", "Timestamp")
	fmt.Println("---------------------------------------------------------------")

	for _, r := range requests {
		fmt.Printf("%-5d %-15s %-8s %-20s %s\n", r.ID, r.IP, r.Method, r.Path, r.Timestamp)
	}
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

	fmt.Println("Connected to Mirage database.")

	store := NewPostgresStore(db)
	run(store)
}
