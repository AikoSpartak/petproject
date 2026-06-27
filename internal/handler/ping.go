package handler

import "net/http"

func Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Pong"))
	if err != nil {
		return
	}

}
