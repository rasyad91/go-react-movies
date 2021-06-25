import React, { Component, Fragment } from 'react';
import './EditMovie.css';
import InputTextArea from './forms/InputTextArea';
import InputText from './forms/InputText'
import InputSelect from './forms/InputSelect'
import Alert from './ui/Alert'
import { Link } from 'react-router-dom';
import { confirmAlert } from 'react-confirm-alert'; // Import
import 'react-confirm-alert/src/react-confirm-alert.css'; // Import css



export default class EditMovie extends Component {

    constructor(props) {
        super(props);
        this.state = {
            movie: {
                id: 0,
                title: "",
                description: "",
                release_date: "",
                runtime: 0,
                rating: 0,
                mpaa_rating: "",
            },
            mpaaOptions: [
                { id: "G", value: "G" },
                { id: "PG13", value: "PG13" },
                { id: "NC17", value: "NC17" },
                { id: "R", value: "R" },
            ],
            isLoaded: false,
            error: null,
            errors: [],
            alert: {
                type: "d-none",
                message: ""
            }
        }

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit = (e) => {
        e.preventDefault();

        //client side validation
        let errors = [];
        if (this.state.movie.title === "") errors.push("title");
        if (this.state.movie.release_date === "") errors.push("release_date");
        if (this.state.movie.runtime === 0) errors.push("runtime");
        if (this.state.movie.mpaa_rating === "") errors.push("mpaa_rating");
        if (this.state.movie.rating === 0) errors.push("rating");
        if (this.state.movie.description === "") errors.push("description");


        this.setState({ errors: errors });
        if (errors.length > 0) return false;

        console.log("Form was submittd");
        const data = new FormData(e.target);
        const payload = Object.fromEntries(data.entries());
        console.log(payload)

        const requestOptions = {
            method: "POST",
            body: JSON.stringify(payload)
        }
        console.log("reqoptions: ", requestOptions)
        fetch("http://localhost:4000/v1/admin/addMovie", requestOptions)
            .then(res => res.json())
            .then(data => {
                if (data.error) { this.setState({ alert: { type: "alert-danger", message: data.error.message } }) }
                if (!data.error) { 
                    this.props.history.push({
                        pathname: "/admin"
                    })}
            })
    }


    hasError = key => { return this.state.errors.indexOf(key) !== -1 }

    handleChange = (e) => {
        let value = e.target.value;
        let name = e.target.name;
        this.setState(prevState => ({
            movie: {
                ...prevState.movie,
                [name]: value
            }
        }))
    }

    confirmDelete = e => {
        console.log("Delete this id" + this.state.movie.id)
        confirmAlert({
            title: 'Delete movie?',
            message: 'Are you sure?',
            buttons: [
                {
                    label: 'Yes',
                    onClick: () => {
                        fetch("http://localhost:4000/v1/admin/deleteMovie/" + this.state.movie.id)
                            .then(res => res.json())
                            .then(data => {
                                if (data.error) { this.setState({ alert: { type: "alert-danger", message: data.error.message } }) }
                                if (!data.error) { 
                                    this.props.history.push({
                                        pathname: "/admin"
                                    })}
                            })
                    }
                },
                {
                    label: 'No',
                    onClick: () => { }
                }
            ]
        });

    }

    componentDidMount() {
        const id = this.props.match.params.id;
        console.log("id: ", id)
        if (id != 0) {
            fetch("http://localhost:4000/v1/movies/" + id)
                .then(res => {
                    if (res.status !== 200) {
                        let err = Error;
                        err.message = "Invalid response code: " + res.status;
                        this.setState({ error: err });
                    }
                    return res.json();
                })
                .then(data => {
                    const releaseDate = new Date(data.movie.release_date);
                    this.setState({
                        movie: {
                            id: data.movie.id,
                            title: data.movie.title,
                            release_date: releaseDate.toISOString().split("T")[0],
                            runtime: data.movie.runtime,
                            mpaa_rating: data.movie.mpaa_rating,
                            rating: data.movie.rating,
                            description: data.movie.description
                        },
                    },
                        error => {
                            this.setState({
                                error
                            });
                        })
                })
        }
        this.setState({ isLoaded: true })
    }

    render() {
        let { movie, isLoaded, error } = this.state
        if (error) return <p>Error: {error.message}</p>
        if (!isLoaded) return <p>Loading...</p>

        return (
            <Fragment>
                <h2>Edit/Add movie</h2>
                <Alert alertType={this.state.alert.type} alertMessage={this.state.alert.message} />


                <form onSubmit={this.handleSubmit}>
                    <input
                        type="hidden"
                        id="id"
                        name="id"
                        value={movie.id}
                        onChange={this.handleChange}
                    ></input>
                    <InputText
                        title="Title"
                        className={this.hasError("title") ? "is-invalid" : ""}
                        errorDiv={this.hasError("title") ? "text-danger" : "d-none"}
                        errorMsg={this.hasError("title") ? "Title required" : ""}
                        type="text"
                        name="title"
                        value={movie.title}
                        handleChange={this.handleChange}
                    />
                    <InputText
                        title="Release Date"
                        className={this.hasError("release_date") ? "is-invalid" : ""}
                        errorDiv={this.hasError("release_date") ? "text-danger" : "d-none"}
                        errorMsg={this.hasError("release_date") ? "Release date required" : ""}
                        type="date"
                        name="release_date"
                        value={movie.release_date}
                        handleChange={this.handleChange}
                    />
                    <InputText
                        title="Runtime"
                        className={this.hasError("runtime") ? "is-invalid" : ""}
                        errorDiv={this.hasError("runtime") ? "text-danger" : "d-none"}
                        errorMsg={this.hasError("runtime") ? "Runtime required" : ""}
                        type="text"
                        name="runtime"
                        value={movie.runtime}
                        handleChange={this.handleChange}
                    />
                    <InputSelect
                        title="MPAA Rating"
                        className={this.hasError("mpaa_rating") ? "is-invalid" : ""}
                        errorDiv={this.hasError("mpaa_rating") ? "text-danger" : "d-none"}
                        errorMsg={this.hasError("mpaa_rating") ? "MPAA rating required" : ""}
                        name="mpaa_rating"
                        options={this.state.mpaaOptions}
                        value={movie.mpaa_rating}
                        handleChange={this.handleChange}
                        placeholder="Choose.."
                    />
                    <InputText
                        title="Rating"
                        className={this.hasError("rating") ? "is-invalid" : ""}
                        errorDiv={this.hasError("rating") ? "text-danger" : "d-none"}
                        errorMsg={this.hasError("rating") ? "Rating required" : ""}
                        type="text"
                        name="rating"
                        value={movie.rating}
                        handleChange={this.handleChange}
                    />
                    <InputTextArea
                        title="Description"
                        className={this.hasError("description") ? "is-invalid" : ""}
                        errorDiv={this.hasError("description") ? "text-danger" : "d-none"}
                        errorMsg={this.hasError("description") ? "Description required" : ""}
                        name="description"
                        value={movie.description}
                        rows="3"
                        handleChange={this.handleChange}
                    />
                    <hr></hr>
                    <button className="btn btn-primary">Save</button>
                    <Link to="/admin" className="btn btn-warning ms-2">Cancel</Link>
                    {movie.id > 0 && (
                        <a onClick={this.confirmDelete} className="btn btn-danger float-end">
                            Delete
                        </a>
                    )}
                </form>

                <div className="mt-3">
                    <pre>{JSON.stringify(this.state, null, 3)}</pre>
                </div>
            </Fragment>
        );
    }
}