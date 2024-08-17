package main

import (
	"log/slog"
	"net/http"
	"os"
)

var links = map[string]string{
	"1": "https://forrostrasbourg.fr/evenements/240821/",
	"2": "https://forrostrasbourg.fr/evenements/240806/",
}

func main() {
	hostPort := os.ExpandEnv("${HOST}:${PORT}")
	if hostPort == ":" {
		hostPort = ":4444"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.EscapedPath()[1:]

		v, ok := links[key]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, v, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/qr/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.EscapedPath()[4:]

		v, ok := links[key]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
		}
		w.Write([]byte("MY QR:" + v))
	})

	slog.Info("listening on", "hostPort", hostPort)
	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		slog.Error("http server interruption", "error", err)
	}
}
