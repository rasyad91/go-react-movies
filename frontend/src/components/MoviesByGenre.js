import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';

export default class Movies extends Component {
    state = {
        movies: [],
        isLoaded: false,
        genreName: "",
        error: null
    };

    componentDidMount() {
        fetch("http://localhost:4000/v1/genres/" + this.props.match.params.id)
            .then(res => {
                if (res.status !== 200) {
                    let err = Error;
                    err.message = "Invalid response code: " + res.status;
                    this.setState({ error: err });
                }
                return res.json();
            })
            .then(data => (
                this.setState({
                    movies: data.movies,
                    isLoaded: true,
                    genreName: this.props.location.genreName
                },
                    error => {
                        this.setState({
                            isLoaded: true,
                            error
                        });
                    })
            ))
    }

    render() {
        let { movies, isLoaded, error, genreName } = this.state
        if (error) return <p>Error: {error.message}</p>
        if (!isLoaded) return <p>Loading...</p>
        if (!movies) { movies = [] }
        return (
            <Fragment>
                <h2>Genre: {genreName}</h2>
                <div className="list-group">
                    {movies.map(m => (
                        <Link key={m.id} className="list-group-item list-group-item-action" to={`/movies/${m.id}`}>{m.title}</Link>
                    ))}
                </div>
            </Fragment>
        );
    }
}