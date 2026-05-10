package main

import (
	"database/sql"
	"log"
)

type Store interface {
	Insert(ip, method, path string)
	GetAll() []Request
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Insert(ip, method, path string) {
	_, err := s.db.Exec(
		"INSERT INTO requests (ip, method, path) VALUES ($1, $2, $3)",
		ip, method, path,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *PostgresStore) GetAll() []Request {
	rows, err := s.db.Query("SELECT id, ip, method, path, timestamp FROM requests ORDER BY id DESC")
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
