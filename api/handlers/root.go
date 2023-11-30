package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Root struct {
	l *log.Logger
}

func NewRoot(l * log.Logger) *Root {
	return &Root{l}
}

func (h *Root) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Serving Root")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Welcome to Chirpy %s\n", d)
}