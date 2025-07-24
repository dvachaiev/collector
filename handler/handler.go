package handler

import (
	"io"
	"log/slog"
	"net/http"
)

func WriteTo(wOut io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(wOut, r.Body); err != nil {
			slog.Warn("Processing request error", "src", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}
