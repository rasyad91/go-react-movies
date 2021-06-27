import React, { Component, Fragment } from 'react';
import InputText from './forms/InputText';

export default class Graphql extends Component {
    constructor(props) {
        super(props);
        this.state = {
            movies: [],
            searchTerm: "",
            isLoaded: false,
            error: null,
            alert: {
                type: "d-none",
                message: "",
            }
        }
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount() {
        const payload = `
        {
            list {
                id
                title
                runtime
                year
                description
            }
        }
        `

        const header = new Headers();
        header.append("Content-Type", "application/json")
        const requestOptions = {
            method: "POST",
            headers: header,
            body: payload,
        }
        fetch(`${process.env.REACT_APP_API_URL}/v1/graphql`, requestOptions)
            .then(res => res.json())
            .then(data => {
                let movieList = Object.values(data.data.list)
                return movieList
            })
            .then(list => {
                this.setState({
                    movies: list
                })
            })
    }

    handleChange = (e) => {
        let value = e.target.value
        console.log("Value: handlechange: ", value)
        this.setState({searchTerm: value}, ()=> console.log("searchterm",this.state.searchTerm ))
        console.log("searchterm", this.state.searchTerm)
        if (value.length > 2) {
            this.performSearch()
        } else {
            this.setState({movies: []})
        }
    }

    performSearch() {
        console.log("payload", this.state.searchTerm)

        const payload = `
        {
            search(titleContains: "${this.state.searchTerm}") {
                id
                title
                runtime
                year
                description
            }
        }
        `
        const header = new Headers();
        header.append("Content-Type", "application/json")
        const requestOptions = {
            method: "POST",
            headers: header,
            body: payload,
        }
        fetch(`${process.env.REACT_APP_API_URL}/v1/graphql`, requestOptions)
            .then(res => res.json())
            .then(data => {
                let movieList = Object.values(data.data.search)
                return movieList
            })
            .then(list => {
                if (list.length > 0) this.setState({ movies: list })
                if (list.length = 0) this.setState({ movies: [] })
            })
    }

    render() {
        let { movies } = this.state
        return (
            <Fragment>
                <h2>Graphql</h2>
                <hr />

                <InputText
                    title="Search"
                    type="text"
                    name="search"
                    value={this.state.searchTerm}
                    handleChange={this.handleChange}
                />
                <div className="list-group">
                    {movies.map(m => (
                        <a
                            key={m.id}
                            className="list-group-item list-group-item-action"
                            href={`/graphql/${m.id}`}
                        >
                            <strong>{m.title}</strong>
                            <small className="text-muted"> ({m.year}) - ({m.runtime}) min</small>
                            <br />
                            {m.description.slice(0, 100)}...
                        </a>
                    ))}
                </div>
            </Fragment>
        )
    }
}