package main

import (
	"backend/internal/handler"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func route() http.Handler {
	router := httprouter.New()
	secure := alice.New(checkToken)

	router.HandlerFunc(http.MethodGet, "/status", handler.Repo.Status)
	router.HandlerFunc(http.MethodPost, "/v1/signin", handler.Repo.Signin)

	router.HandlerFunc(http.MethodGet, "/v1/movies", handler.Repo.GetAllMovies)
	// :id because of julienschmidt package
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", handler.Repo.GetMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", handler.Repo.GetAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/genres/:id", handler.Repo.GetAllMoviesByGenre)

	router.POST("/v1/admin/editMovie", wrap(secure.ThenFunc(handler.Repo.AddMovie)))
	router.HandlerFunc(http.MethodPost, "/v1/admin/addMovie", handler.Repo.AddMovie)

	router.HandlerFunc(http.MethodGet, "/v1/admin/deleteMovie/:id", handler.Repo.DeleteMovie)

	return enableCORS(router)
}
