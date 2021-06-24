package postgres

import (
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// dbRepo wrapper for database
type dbRepo struct {
	*sql.DB
}

// New returns Database repo with db pool
func New(conn *sql.DB) repository.Database {
	return &dbRepo{
		DB: conn,
	}
}

// GetMovideByID returns one movie and error
func (m *dbRepo) GetMovieByID(id int) (model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, title, description, year, release_date, runtime, mpaa_rating, created_at, updated_at
			  FROM movies
			  WHERE id = $1`

	movie := model.Movie{}
	if err := m.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	); err != nil {
		return model.Movie{}, err
	}
	movie.Genre, _ = m.GetGenreByID(movie.ID)
	fmt.Println("movie:", movie)

	return movie, nil
}

// GetAllMovies returns a slice of models and error
func (m *dbRepo) GetAllMovies() ([]model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, title, description, year, release_date, runtime, mpaa_rating, created_at, updated_at
			  FROM movies 
			  ORDER BY title`
	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		movie := model.Movie{}
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		); err != nil {
			return nil, nil
		}
		movie.Genre, _ = m.GetGenreByID(movie.ID)
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, nil
	}
	return movies, nil
}

// GetGenreByMovieID returns the genres in the movie
func (m *dbRepo) GetGenreByID(id int) (map[int]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT g.id, g.genre_name, g.created_at, g.updated_at
	FROM movies_genres mg
	LEFT JOIN genres g on (g.id = mg.genre_id)
	WHERE mg.movie_id = $1`

	rows, err := m.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make(map[int]string)
	for rows.Next() {
		g := model.Genre{}
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		genres[g.ID] = g.Name
	}
	return genres, nil
}
