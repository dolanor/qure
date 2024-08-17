package main

import "net/http"

func qrCodeHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.EscapedPath()[4:]

	v, ok := links[key]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
	}
	w.Write([]byte("MY QR:" + v))
}
