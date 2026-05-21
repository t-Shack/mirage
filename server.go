package main

import (
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
	fmt.Fprintln(w, "<html><body><h1>Welcome</h1></body></html>")
}

func (s *Server) Start(port string) {
	http.HandleFunc("/", s.handleRequest)

	log.Printf("Mirage listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
