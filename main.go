package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

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

	createTable := `
	CREATE TABLE IF NOT EXISTS requests (
		id        SERIAL PRIMARY KEY,
		ip        TEXT NOT NULL,
		method    TEXT NOT NULL,
		path      TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT NOW()
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table ready.")

	_, err = db.Exec(
		"INSERT INTO requests (ip, method, path) VALUES ($1, $2, $3)",
		"192.168.1.1", "GET", "/admin",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Row inserted.")

	var ip, method, path string
	row := db.QueryRow("SELECT ip, method, path FROM requests ORDER BY id DESC LIMIT 1")
	err = row.Scan(&ip, &method, &path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last request -> IP: %s | Method: %s | Path: %s\n", ip, method, path)
}
