package api

import (
	"fmt"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("route %s not found", r.URL.Path), http.StatusNotFound)
}
