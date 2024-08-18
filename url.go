package main

import (
	"net/http"
)

func urlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.EscapedPath()[1:]

	v, ok := links[key]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, v, http.StatusTemporaryRedirect)
}
