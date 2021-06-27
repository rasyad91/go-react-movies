import React, { Component, Fragment } from 'react'

export default class ShowOneMovie extends Component {
    state = {
        movie: {},
        isLoaded: false,
        error: null
    };

    componentDidMount() {
        fetch(`${process.env.REACT_APP_API_URL}/v1/movies/` + this.props.match.params.id)
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
                    movie: data.movie,
                    isLoaded: true
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
        const { movie, isLoaded, error } = this.state
        if (error) return <p>Error: {error.message}</p>
        if (!isLoaded) return <p>Loading...</p>
        if (movie.genres) {
            movie.genres = Object.values(movie.genres)
        };
        if (!movie.genres) { movie.genres = [] }

        return (
            <Fragment>
                <h2>Movie: {movie.title} ({movie.year})</h2>
                {movie.poster !== "" && (
                    <div>
                        <img src={`https://image.tmdb.org/t/p/w200${movie.poster}`} alt="poster"></img>
                    </div>
                )}
                <div className="float-start">
                    <small>Rating: {movie.mpaa_rating}</small>
                </div>
                <div className="float-end">
                    {movie.genres.map((m) => (
                        <span key={m} className="badge bg-secondary me-1">
                            {m}
                        </span>
                    ))}
                </div>
                <div className="clearfix"></div>
                <hr></hr>
                <table className="table table-compact table-light">
                    <thead></thead>
                    <tbody>
                        <tr>
                            <td><strong>Title:</strong></td>
                            <td><strong>{movie.title}</strong></td>
                        </tr>
                        <tr>
                            <td><strong>Year:</strong></td>
                            <td><strong>{movie.year}</strong></td>
                        </tr>
                        <tr>
                            <td><strong>Description:</strong></td>
                            <td><strong>{movie.description}</strong></td>
                        </tr>
                        <tr>
                            <td><strong>Runtime:</strong></td>
                            <td><strong>{movie.runtime} min</strong></td>
                        </tr>
                        <tr>
                            <td><strong>Rating:</strong></td>
                            <td><strong>{movie.rating}</strong></td>
                        </tr>
                    </tbody>
                </table>
            </Fragment>
        )
    }
}