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

func insertRequest(db *sql.DB, ip, method, path string) {
	_, err := db.Exec(
		"INSERT INTO requests (ip, method, path) VALUES ($1, $2, $3)",
		ip, method, path,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRequests(db *sql.DB) []Request {
	rows, err := db.Query("SELECT id, ip, method, path, timestamp FROM requests ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var results []Request

	for rows.Next() {
		var r Request
		err := rows.Scan(&r.ID, &r.IP, &r.Method, &r.Path, &r.Timestamp)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, r)
	}

	return results
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

	insertRequest(db, "10.0.0.1", "POST", "/login")
	insertRequest(db, "172.16.0.5", "GET", "/etc/passwd")
	insertRequest(db, "192.168.1.100", "GET", "/wp-admin")

	fmt.Println("Inserted 3 requests.")

	requests := getAllRequests(db)

	fmt.Printf("\n%-5s %-15s %-8s %-20s %s\n", "ID", "IP", "Method", "Path", "Timestamp")
	fmt.Println("---------------------------------------------------------------")

	for _, r := range requests {
		fmt.Printf("%-5d %-15s %-8s %-20s %s\n", r.ID, r.IP, r.Method, r.Path, r.Timestamp)
	}
}
