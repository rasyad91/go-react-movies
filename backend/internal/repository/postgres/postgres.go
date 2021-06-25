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

	query := `SELECT id, title, description, year, release_date, runtime, mpaa_rating, rating, created_at, updated_at
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
		&movie.Rating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	); err != nil {
		return model.Movie{}, err
	}
	movie.Genres, _ = m.GetGenreByMovieID(movie.ID)
	return movie, nil
}

// GetAllMovies returns a slice of models and error, insert int if want to filter by genre
func (m *dbRepo) GetAllMovies(genre ...int) ([]model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("WHERE id in (SELECT movie_id FROM movies_genres WHERE genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`SELECT id, title, description, year, release_date, runtime, mpaa_rating, rating, created_at, updated_at
			  FROM movies
			  %s 
			  ORDER BY title`, where)

	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return nil, err
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
			&movie.Rating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		); err != nil {
			return nil, err
		}
		movie.Genres, _ = m.GetGenreByMovieID(movie.ID)
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

// GetGenreByMovieID returns the genres in the movie
func (m *dbRepo) GetGenreByMovieID(id int) (map[int]string, error) {
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

// GetAllGenres returns all genres and error
func (m *dbRepo) GetAllGenres() ([]model.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, genre_name, created_at, updated_at
			  FROM genres`

	genres := []model.Genre{}

	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		genres = append(genres, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return genres, nil
}

func (m *dbRepo) InsertMovie(movie model.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO 
			  movies (title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  `
	if _, err := m.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		time.Now(),
		time.Now(),
	); err != nil {
		return err
	}

	return nil
}

func (m *dbRepo) UpdateMovie(movie model.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE movies
			  SET 	title = $1, 
			  		description = $2, 
					year = $3, 
					release_date = $4, 
					runtime = $5, 
					rating = $6, 
					mpaa_rating = $7,
					updated_at = $8
			  WHERE id = $9
			  `
	if _, err := m.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		time.Now(),
		movie.ID,
	); err != nil {
		return err
	}

	return nil
}

func (m *dbRepo) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM movies
			  WHERE id = $1
			  `
	if _, err := m.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
