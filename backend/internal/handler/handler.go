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
	"net/url"
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
	fmt.Println("Getmovie params:", params)

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

	if movie.Poster == "" {
		fmt.Println("Get movie Poster")
		movie = getPoster(movie)
		fmt.Println("Return from movie Poster")

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

	id, _ := strconv.Atoi(payload.ID)

	if id != 0 {
		m, _ := m.DB.GetMovieByID(id)
		movie = m
	}

	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.MPAARating = payload.MPAARating
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.Rating, _ = strconv.Atoi(payload.Rating)

	if movie.Poster == "" {
		movie = getPoster(movie)
	}

	log.Println(movie)
	if movie.ID == 0 {
		if err := m.DB.InsertMovie(movie); err != nil {
			fmt.Println(err)
			util.ErrorJSON(w, err)
			return
		}
	} else {
		if err := m.DB.UpdateMovie(movie); err != nil {
			fmt.Println(err)
			util.ErrorJSON(w, err)
			return
		}
	}

	resp := jsonResp{OK: true}
	if err := util.WriteJSON(w, "response", resp); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

}

func (m *Repository) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	// params := httprouter.ParamsFromContext(r.Context())
	params := r.Context().Value("params").(httprouter.Params)
	fmt.Println("Deletemovie ctx:", r.Context())

	fmt.Println("deletemovie params:", params)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		m.App.Logger.Println(errors.New("invalid id parameter"))
		util.ErrorJSON(w, err)
		return
	}
	if err := m.DB.DeleteMovie(id); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

	resp := jsonResp{OK: true}
	if err := util.WriteJSON(w, "response", resp); err != nil {
		m.App.Logger.Println(err)

	}
}

func getPoster(movie model.Movie) model.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}
	client := &http.Client{}
	domain := "https://api.themoviedb.org/3/search/movie?"
	apiKey := "83cc89606f90461a781921ec33c57c57"
	query := "&query="
	url := fmt.Sprintf("%sapi_key=%s%s%s", domain, apiKey, query, url.QueryEscape(movie.Title))
	fmt.Println("NEW REQUEST")
	fmt.Println("NEW REQUEST")
	fmt.Println("NEW REQUEST")

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return movie
	}
	fmt.Println("RETURN NEW REQUEST")
	fmt.Println("RETURN NEW REQUEST")
	fmt.Println("RETURN NEW REQUEST")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	response, err := client.Do(r)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer response.Body.Close()

	fmt.Println("RETURN RESPOSE")
	fmt.Println("RETURN RESPOSE")
	fmt.Println("RETURN RESPOSE")
	fmt.Printf("response %#v", response)

	// bodyBytes, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return movie
	// }
	responseObject := &TheMovieDB{}
	enc := json.NewDecoder(response.Body)
	enc.Decode(responseObject)
	fmt.Println("DECODE")
	fmt.Println("DECODE")
	fmt.Println("DECODE")

	if len(responseObject.Results) > 0 {
		movie.Poster = responseObject.Results[0].PosterPath
		fmt.Println("Poster PATH: ", movie.Poster)
	}

	return movie
}
