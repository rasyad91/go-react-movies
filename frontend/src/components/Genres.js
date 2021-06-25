import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom'

export default class Genres extends Component {

    state = {
        genres: [],
        isLoaded: false,
        error: null
    }

    componentDidMount() {
        fetch("http://localhost:4000/v1/genres")
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
                    genres: data.genres,
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
        const { genres, isLoaded, error } = this.state
        if (error) return <p>Error: {error.message}</p>
        if (!isLoaded) return <p>Loading...</p>

        return (
            <Fragment>
                <h2>Genres</h2>
                <div className="list-group">
                    {genres.map(g => (
                        <Link
                            key={g.id}
                            className="list-group-item list-group-item-active"
                            to={{
                                pathname: `/genres/${g.id}`,
                                genreName: g.name
                            }}>
                            {g.name}
                        </Link>
                    ))}

                </div>
            </Fragment>
        );
    }
}