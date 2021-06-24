package repository

import "backend/internal/model"

type Database interface {
	GetMovieByID(id int) (model.Movie, error)
	GetAllMovies() ([]model.Movie, error)
}
