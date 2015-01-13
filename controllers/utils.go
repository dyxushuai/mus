package controllers


import (
	"net/http"
	"fmt"
)

func JsonView(fn func(w http.ResponseWriter, r *http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data := fn(w, r)
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "application/json", "utf-8"))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		fmt.Fprintf(w, data)
	})
}
