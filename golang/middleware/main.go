package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Middleware func(http.Handler) http.Handler

func httpBasicAuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(fmt.Sprintf("failed to get credentials")))
			return
		}
		log.Printf("decoded basic auth: %s:%s", user, pass)

		handler.ServeHTTP(rw, req)
	})
}

func loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		handler.ServeHTTP(rw, r)
	})
}

func main() {
	fsys := os.DirFS(".")
	handler := http.FileServerFS(fsys)

	middlewares := []Middleware{
		httpBasicAuthMiddleware,
		loggingMiddleware,
	}
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	server := http.Server{
		Addr:    "localhost:9772",
		Handler: handler,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}
