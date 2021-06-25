package main

import (
	"backend/internal/handler"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func route() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", handler.Repo.Status)

	router.HandlerFunc(http.MethodGet, "/v1/movies", handler.Repo.GetAllMovies)
	// :id because of julienschmidt package
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", handler.Repo.GetMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", handler.Repo.GetAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/genres/:id", handler.Repo.GetAllMoviesByGenre)

	router.HandlerFunc(http.MethodPost, "/v1/admin/addMovie", handler.Repo.AddMovie)

	return enableCORS(router)
}
