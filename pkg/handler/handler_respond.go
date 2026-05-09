package handler

import (
	"net/http"
)

func Respond(body string, status int, headers map[string]string) http.Handler {
	if status == 0 {
		status = http.StatusOK
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Set(k, v)
		}

		w.WriteHeader(status)
		w.Write([]byte(body))
	})
}
