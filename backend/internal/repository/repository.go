package repository

import "backend/internal/model"

type Database interface {
	GetMovieByID(id int) (model.Movie, error)
	GetAllMovies(genre ...int) ([]model.Movie, error)
	InsertMovie(movie model.Movie) error
	UpdateMovie(movie model.Movie) error
	DeleteMovie(id int) error

	GetAllGenres() ([]model.Genre, error)
}
