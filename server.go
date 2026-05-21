package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Server struct {
	store Store
}

func NewServer(store Store) *Server {
	return &Server{store: store}
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	method := r.Method
	path := r.URL.Path

	s.store.Insert(ip, method, path)

	log.Printf("Captured -> IP: %s | Method: %s | Path: %s", ip, method, path)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "<html><body><h1>Welcome to Mirage!</h1></body></html>")
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	requests, err := s.store.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requests)
	if err != nil {
		log.Printf("Failed to encode response: %s", err)
	}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/admin", s.handleAdmin)
	http.HandleFunc("/", s.handleRequest)

	log.Printf("Mirage listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
