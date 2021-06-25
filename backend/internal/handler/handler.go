package handler

import (
	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/util"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Repository struct {
	App *config.Application
	DB  repository.Database
}

var Repo *Repository

func Initialize(app *config.Application, db repository.Database) *Repository {
	return &Repository{
		App: app,
		DB:  db,
	}
}

func New(r *Repository) {
	Repo = r
}

func (m *Repository) Status(w http.ResponseWriter, r *http.Request) {
	currentStatus := model.AppStatus{
		Status:      "Available",
		Environment: m.App.Env,
		Version:     m.App.Version,
	}

	b := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b.SetIndent("", "\t")
	if err := b.Encode(currentStatus); err != nil {
		m.App.Logger.Println(err)
	}
}

func (m *Repository) GetMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		m.App.Logger.Println(errors.New("invalid id parameter"))
		util.ErrorJSON(w, err)
		return
	}

	movie, err := m.DB.GetMovieByID(id)
	if err == sql.ErrNoRows {
		m.App.Logger.Println(errors.New("no movie found"))
		util.ErrorJSON(w, err)
		return
	}
	if err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	if err := util.WriteJSON(w, "movie", movie); err != nil {
		m.App.Logger.Println(err)
	}

}

func (m *Repository) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := m.DB.GetAllMovies()
	if err == sql.ErrNoRows {
		m.App.Logger.Println(errors.New("no movie found"))
		util.ErrorJSON(w, err)
		return
	}
	if err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	if err := util.WriteJSON(w, "movies", movies); err != nil {
		m.App.Logger.Println(err)
	}
}

func (m *Repository) GetAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		m.App.Logger.Println(errors.New("invalid id parameter"))
		util.ErrorJSON(w, err)
		return
	}
	movies, err := m.DB.GetAllMovies(id)
	if err == sql.ErrNoRows {
		m.App.Logger.Println(errors.New("no movie found"))
		util.ErrorJSON(w, err)
		return
	}
	if err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	if err := util.WriteJSON(w, "movies", movies); err != nil {
		m.App.Logger.Println(err)
	}
}

func (m *Repository) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := m.DB.GetAllGenres()
	if err == sql.ErrNoRows {
		m.App.Logger.Println(errors.New("no genres found"))
		util.ErrorJSON(w, err)
		return
	}
	if err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	if err := util.WriteJSON(w, "genres", genres); err != nil {
		m.App.Logger.Println(err)
	}
}

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AddMovie(w http.ResponseWriter, r *http.Request) {
	var payload model.MoviePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	var movie model.Movie
	movie.ID, _ = strconv.Atoi(payload.ID)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.MPAARating = payload.MPAARating
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.Rating, _ = strconv.Atoi(payload.Rating)

	log.Println(movie)
	if err := m.DB.InsertMovie(movie); err != nil {
		fmt.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	resp := jsonResp{OK: true}
	if err := util.WriteJSON(w, "response", resp); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

}

func (m *Repository) DeleteMovie(w http.ResponseWriter, r *http.Request)  {}
func (m *Repository) SearchMovies(w http.ResponseWriter, r *http.Request) {}
