package api

import (
	"net/http"

	"github.com/aljabri00056/installer/handler"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	c := handler.GetConfigFromEnv()

	h := &handler.Handler{Config: c}

	h.ServeHTTP(w, r)
}
