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

	http.HandleFunc("/", urlShortenerHandler)

	http.HandleFunc("/qr/", qrCodeHandler)

	slog.Info("listening on", "hostPort", hostPort)
	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		slog.Error("http server interruption", "error", err)
	}
}
