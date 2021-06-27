import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';

export default class Admin extends Component {
    state = {
        movies: [],
        isLoaded: false,
        error: null
    };

    componentDidMount() {
        if (this.props.jwt === "") return this.props.history.push({pathname: "/"})


        fetch(`${process.env.REACT_APP_API_URL}/v1/movies`)
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
        const { movies, isLoaded, error } = this.state
        if (error) return <p>Error: {error.message}</p>
        if (!isLoaded) return <p>Loading...</p>
        return (
            <Fragment>
                <h2>Manage movies</h2>
                <div className="list-group">
                    {movies.map(m => (
                        <Link key={m.id} className="list-group-item list-group-item-action" to={`/admin/movies/${m.id}`}>{m.title}</Link>
                    ))}
                </div>
            </Fragment>
        );
    }
}