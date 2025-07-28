package handler

import (
	"io"
	"log/slog"
	"net/http"
)

func WriteTo(wOut io.Writer, fn func([]byte) ([]byte, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Warn("Reading request body", "src", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		out, err := fn(body)
		if err != nil {
			slog.Warn("Processing request body", "src", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if _, err = wOut.Write(out); err != nil {
			slog.Warn("Storing processed data", "src", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
