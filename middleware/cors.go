package middleware

import (
	"actionflow/core"
	"net/http"
)

//CORSMiddleware 支持跨域
func CORSMiddleware(srv *core.Server, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, HEAD, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type,authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "false")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
