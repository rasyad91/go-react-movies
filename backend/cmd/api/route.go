package main

import (
	"backend/internal/handler"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func route() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", handler.Repo.Status)

	return router
}
