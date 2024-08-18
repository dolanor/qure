package main

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/dolanor/rip"
	"github.com/dolanor/rip/encoding/html"
	"github.com/gorilla/handlers"
)

func main() {
	cfg := loadConfig()

	ro := rip.NewRouteOptions().
		WithCodecs(
			html.NewEntityCodec("/admin/links/"),
			html.NewEntityFormCodec("/admin/links/"),
		).
		WithErrors(rip.StatusMap{
			ErrEmptyLinkSlug: http.StatusBadRequest,
		}).
		WithMiddlewares(loggerMiddleware(os.Stdout))

	db, err := sql.Open("sqlite", filepath.Join(cfg.dbDirPath, "qure.db"))
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

	http.HandleFunc("/qr/", qrCodeHandler(cfg.domain))

	http.HandleFunc(rip.HandleEntities("/admin/links/", lp, ro))

	slog.Info("listening on", "hostPort", cfg.hostPort)
	err = http.ListenAndServe(cfg.hostPort, nil)
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
