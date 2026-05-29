package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	store Store
}

type AdminPageData struct {
	Requests []Request
	Stats    Stats
	IP       string
	Method   string
	Path     string
	Page     int
	PrevPage int
	NextPage int
	HasNext  bool
	HasPrev  bool
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
	ip := r.URL.Query().Get("ip")
	method := r.URL.Query().Get("method")
	path := r.URL.Query().Get("path")

	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	const limit = 25
	offset := (page - 1) * limit

	filter := Filter{
		IP:     ip,
		Method: method,
		Path:   path,
		Limit:  limit + 1,
		Offset: offset,
	}

	requests, err := s.store.GetAll(filter)
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	hasNext := len(requests) > limit
	if hasNext {
		requests = requests[:limit]
	}

	for i := range requests {
		requests[i].Severity = classifyPath(requests[i].Path)
	}

	stats, err := s.store.GetStats()
	if err != nil {
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template parse error: %s", err)
		return
	}

	data := AdminPageData{
		Requests: requests,
		Stats:    stats,
		IP:       ip,
		Method:   method,
		Path:     path,
		Page:     page,
		PrevPage: page - 1,
		NextPage: page + 1,
		HasNext:  hasNext,
		HasPrev:  page > 1,
	}

	err = tmpl.Execute(w, data)
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

func (s *Server) handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func (s *Server) Start(port string) {
	http.HandleFunc("/static/favicon.ico", s.handleFavicon)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/admin", s.handleAdmin)
	http.HandleFunc("/", s.handleRequest)

	log.Printf("Mirage listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
