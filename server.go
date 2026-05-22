package main

import (
	_ "encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
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

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	requests, err := s.store.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	for i := range requests {
		requests[i].Severity = classifyPath(requests[i].Path)
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template parse error: %s", err)
		return
	}

	err = tmpl.Execute(w, requests)
	if err != nil {
		log.Printf("Template execute error: %s", err)
	}
}

func classifyPath(path string) string {
	dangerous := []string{
		"/etc/passwd", "/etc/shadow", "/.env",
		"/wp-config", "/.git", "/config",
	}
	suspicious := []string{
		"/admin", "/login", "/wp-admin",
		"/phpmyadmin", "/secret", "/backup",
	}

	for _, d := range dangerous {
		if strings.Contains(path, d) {
			return "danger"
		}
	}
	for _, s := range suspicious {
		if strings.Contains(path, s) {
			return "warn"
		}
	}
	return "normal"
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
