package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Healthz struct {
	l *log.Logger
}

func NewHealthz(l * log.Logger) *Healthz {
	return &Healthz{l}
}

func (h *Healthz) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Serving Healthz")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(rw, "%s\n %s\n", http.StatusText(http.StatusOK), d)
}