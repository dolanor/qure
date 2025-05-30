package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/skip2/go-qrcode"
)

func qrCodeHandler(domain string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.EscapedPath()[4:]

		shortURL := fmt.Sprintf("http://%s/%s", domain, key)

		uppercaseURL := strings.ToUpper(shortURL)

		png, err := qrcode.Encode(uppercaseURL, qrcode.Low, 256)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "image/png")

		_, err = w.Write(png)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
