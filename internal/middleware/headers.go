package middleware

import "net/http"

func MiddlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(rw, r)
	})
}