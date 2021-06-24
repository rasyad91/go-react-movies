package handler

import (
	"backend/internal/config"
	"backend/internal/model"
	"encoding/json"
	"net/http"
)

type Repository struct {
	*config.Application
}

var Repo *Repository

func Initialize(app *config.Application) *Repository {
	return &Repository{
		app,
	}
}

func New(r *Repository) {
	Repo = r
}

func (m *Repository) Status(w http.ResponseWriter, r *http.Request) {
	currentStatus := model.AppStatus{
		Status:      "Available",
		Environment: m.Env,
		Version:     m.Version,
	}

	b := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b.SetIndent("", "\t")
	if err := b.Encode(currentStatus); err != nil {
		m.Logger.Println(err)
	}

	// b, err := json.MarshalIndent(currentStatus, "", "\t")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Fprint(w, b)

}
