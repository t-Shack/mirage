package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Store interface {
	Insert(ip, method, path string)
	GetAll(filter Filter) ([]Request, error)
	GetStats() (Stats, error)
}

type Filter struct {
	IP     string
	Method string
	Path   string
	Limit  int
	Offset int
}

type Stats struct {
	TotalProbes int
	Last24h     int
	UniqueIPs   int
	TopPath     string
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
		fmt.Println("Insert error:", err)
	}
}

func (s *PostgresStore) GetAll(filter Filter) ([]Request, error) {
	query := "SELECT id, ip, method, path, timestamp FROM requests WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if filter.IP != "" {
		query += fmt.Sprintf(" AND ip ILIKE $%d", idx)
		args = append(args, "%"+filter.IP+"%")
		idx++
	}
	if filter.Method != "" {
		query += fmt.Sprintf(" AND method = $%d", idx)
		args = append(args, filter.Method)
		idx++
	}
	if filter.Path != "" {
		query += fmt.Sprintf(" AND path ILIKE $%d", idx)
		args = append(args, "%"+filter.Path+"%")
		idx++
	}

	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Request

	for rows.Next() {
		var r Request
		var ts time.Time
		err := rows.Scan(&r.ID, &r.IP, &r.Method, &r.Path, &ts)
		if err != nil {
			return nil, err
		}
		r.Timestamp = ts.Format("2006-01-02 15:04:05")
		results = append(results, r)
	}

	return results, nil
}

func (s *PostgresStore) GetStats() (Stats, error) {
	var stats Stats

	err := s.db.QueryRow("SELECT COUNT(*) FROM requests").Scan(&stats.TotalProbes)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM requests WHERE timestamp > NOW() - INTERVAL '24 hours'",
	).Scan(&stats.Last24h)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow(
		"SELECT COUNT(DISTINCT ip) FROM requests",
	).Scan(&stats.UniqueIPs)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow(
		"SELECT path FROM requests GROUP BY path ORDER BY COUNT(*) DESC LIMIT 1",
	).Scan(&stats.TopPath)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	return stats, nil
}
