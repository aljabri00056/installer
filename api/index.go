package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/divyam234/installer/handler"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	c := handler.DefaultConfig

	if c.Token == "" && os.Getenv("GH_TOKEN") != "" {
		c.Token = os.Getenv("GH_TOKEN")
	}
	if c.ForceUser != "" {
		log.Printf("locked user to '%s'", c.ForceUser)
	}
	if c.ForceRepo != "" {
		log.Printf("locked repo to '%s'", c.ForceRepo)
	}

	h := &handler.Handler{Config: c}

	h.ServeHTTP(w, r)
}
