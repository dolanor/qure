package main

import (
	"log/slog"
	"net/http"
)

func urlShortenerHandler(lp *ShortURLProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.EscapedPath()[1:]
		slog.Info("click", "slug", slug)

		l, err := lp.FindBySlug(r.Context(), slug)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, l.URL, http.StatusTemporaryRedirect)

		err = lp.Click(r.Context(), slug)
		if err != nil {
			http.Error(w, "can not update click counter", http.StatusInternalServerError)
			return
		}

	}
}
