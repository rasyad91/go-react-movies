package handler

import (
	"backend/internal/model"
	"backend/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

var movies []model.Movie

// populate fields with database fields
var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:       "Movie",
		Interfaces: nil,
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"poster": &graphql.Field{
				Type: graphql.String,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
		Description: "",
	},
)

//graphql scehma definition
var fields = graphql.Fields{
	"movie": &graphql.Field{
		Type: movieType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: nil,
				Description:  "movie id",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			if ok {
				for _, movie := range movies {
					if movie.ID == id {
						return movie, nil
					}
				}
			}
			return nil, nil
		},
		Description: "Get movie by id",
	},
	"list": &graphql.Field{
		Type: graphql.NewList(movieType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			fmt.Println("list RESOLVE FUNC:")
			return movies, nil
		},
		Description: "Get all movies",
	},
	"search": &graphql.Field{
		Type: graphql.NewList(movieType),
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var list []model.Movie
			search, ok := p.Args["titleContains"].(string)
			if ok {
				for _, v := range movies {
					if strings.Contains(v.Title, search) {
						log.Println("found one")
						list = append(list, v)
					}
				}
			}
			return list, nil
		},
		Description: "Search movie by title",
	},
}

func (m *Repository) ListMovies(w http.ResponseWriter, r *http.Request) {
	movies, _ = m.DB.GetAllMovies()

	q, _ := io.ReadAll(r.Body)
	query := string(q)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, errors.New("fail to create schema"))
		return
	}
	params := graphql.Params{Schema: schema, RequestString: query}
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		m.App.Logger.Println(result.Errors)
		util.ErrorJSON(w, fmt.Errorf("fail to execute graphql operation: %#v", result.Errors))
		return
	}
	fmt.Printf("%#v\n", result)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(result)

	// if err := util.WriteJSON(w, "result", result); err != nil {
	// 	m.App.Logger.Println(err)
	// 	util.ErrorJSON(w, err)
	// }
}
