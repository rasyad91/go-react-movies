package handler

import (
	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/util"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	if err := util.WriteJSON(w, movie); err != nil {
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

	if err := util.WriteJSON(w, movies); err != nil {
		m.App.Logger.Println(err)
	}
}

func (m *Repository) DeleteMovie(w http.ResponseWriter, r *http.Request)  {}
func (m *Repository) InsertMovie(w http.ResponseWriter, r *http.Request)  {}
func (m *Repository) UpdateMovie(w http.ResponseWriter, r *http.Request)  {}
func (m *Repository) SearchMovies(w http.ResponseWriter, r *http.Request) {}
