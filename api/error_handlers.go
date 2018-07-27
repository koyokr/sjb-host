package api

import (
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}
