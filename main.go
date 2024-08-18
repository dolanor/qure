package main

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"github.com/dolanor/rip"
	"github.com/dolanor/rip/encoding/html"
	"github.com/gorilla/handlers"
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

	domain := os.Getenv("DOMAIN")

	ro := rip.NewRouteOptions().
		WithCodecs(
			html.NewEntityCodec("/admin/links/"),
			html.NewEntityFormCodec("/admin/links/"),
		).
		WithErrors(rip.StatusMap{
			ErrEmptyLinkSlug: http.StatusBadRequest,
		}).
		WithMiddlewares(loggerMiddleware(os.Stdout))

	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		slog.Error("sqlite open", "error", err)
		os.Exit(1)
	}

	lp, err := NewLinkProvider(db)
	if err != nil {
		slog.Error("url provider creation", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/", urlShortenerHandler)

	http.HandleFunc("/qr/", qrCodeHandler(domain))

	http.HandleFunc(rip.HandleEntities("/admin/links/", lp, ro))

	slog.Info("listening on", "hostPort", hostPort)
	err = http.ListenAndServe(hostPort, nil)
	if err != nil {
		slog.Error("http server listen", "error", err)
		os.Exit(1)
	}
}

func loggerMiddleware(logOut io.Writer) func(http.HandlerFunc) http.HandlerFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		logHandler := handlers.LoggingHandler(logOut, hf)
		return logHandler.ServeHTTP
	}
}
