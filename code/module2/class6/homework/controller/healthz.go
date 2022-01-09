package controller

import "net/http"

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
