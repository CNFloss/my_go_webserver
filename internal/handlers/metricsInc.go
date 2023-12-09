package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sync"
)

type ApiConfig struct {
	l              *log.Logger
	fileserverHits int
	mu             sync.Mutex  // Mutex for safe counter increment
}

type PageData struct {
    Hits int
}

func NewApiConfig(l * log.Logger) *ApiConfig {
	return &ApiConfig{l, 0, sync.Mutex{}}
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cfg.l.Println("incrementing hits")

		cfg.mu.Lock()          // Lock the mutex before modifying the counter
		cfg.fileserverHits += 1
		cfg.mu.Unlock()        // Unlock the mutex after modifying the counter

		next.ServeHTTP(rw, r)
	})
}

func (cfg *ApiConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	cfg.l.Println("Serving Metrics")
	tmpl, err := template.ParseFiles("admin.html")
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	_, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
	
	rw.Header().Add("Content-Type", "text/html; charset=utf-8")
	data := PageData{
		Hits: cfg.fileserverHits,
  }

	err = tmpl.Execute(rw, data)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
}

func (cfg *ApiConfig) ResetHits(rw http.ResponseWriter, r *http.Request) {
	cfg.l.Println("Reseting Metrics")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	cfg.mu.Lock()          // Lock the mutex before modifying the counter
	cfg.fileserverHits = 0
	cfg.mu.Unlock()     
	
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(rw, "Reset Hits: %d\n %s\n", cfg.fileserverHits, d)
}