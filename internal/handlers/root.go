package handlers

import (
	"log"
	"net/http"
	"os"
)

type Root struct {
	l *log.Logger
}

func NewRoot(l * log.Logger) *Root {
	return &Root{l}
}

func (h *Root) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Serving Root")

	htmlFilePath := "index.html"

	data, err := os.ReadFile(htmlFilePath)
	if err != nil {
		h.l.Printf("Error reading HTML file: %v", err)
		http.Error(rw, "Unable to read HTML file", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html")

	_, err = rw.Write(data)
	if err != nil {
		h.l.Printf("Error writing HTML content: %v", err)
		http.Error(rw, "Error serving HTML content", http.StatusInternalServerError)
	}
}