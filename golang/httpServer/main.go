package main

import (
	"embed"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

//go:embed templates/*
var rawTemplates embed.FS

var tpls *template.Template

func init() {
	tpls = template.Must(template.ParseFS(rawTemplates, "templates/*"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", getIndex)
	router.HandleFunc("POST /postForm", postForm)
	http.ListenAndServe("127.0.0.1:8084", logRequestHandler(router))
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	tpls.Execute(w, nil)
}

func postForm(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	slog.Info(string(data))
	header := w.Header()
	header["Location"] = []string{"/"}
	w.WriteHeader(302)
}
func logRequestHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		method := r.Method
		uri := r.URL.String()
		slog.Info("mymsg", "method", method, "uri", uri)
	})
}
