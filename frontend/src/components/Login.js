import React, { Component, Fragment } from 'react'
import InputText from './forms/InputText';
import Alert from './ui/Alert';

export default class Login extends Component {
    constructor(props) {
        super(props);

        this.state = {
            email: "",
            password: "",
            error: null,
            errors: [],
            alert: {
                type: "d-none",
                message: ""
            }
        }
        this.handleChange = this.handleChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)

    }

    handleChange = e => {
        let value = e.target.value
        let name = e.target.name
        this.setState(prev => ({
            ...prev,
            [name]: value
        }))
    }

    handleSubmit = e => {
        e.preventDefault();
        console.log("submitted")
        let errors = []

        if (this.state.email === "") errors.push("email");
        if (this.state.password === "") errors.push("password");

        this.setState({errors: errors})

        if (errors.length !== 0) return false; 

        const formData = new FormData(e.target)
        const payload = Object.fromEntries(formData.entries());
        const request = {
            method: "POST",
            body: JSON.stringify(payload)
        }
        fetch("http://localhost:4000/v1/signin", request)
            .then(res => res.json())
            .then(data => {
                if (data.error) {
                    return this.setState({alert:{type:"alert-danger", message: data.error.message}})
                }
                if (!data.error){
                    this.handleJWTChange(Object.values(data)[0])
                    this.setState({alert:{type:"d-none", message: ""}})
                    this.props.history.push({
                        pathname: "/admin"
                    })
                    console.log(data)
                }
            })
    }

    hasErrors = key => { return this.state.errors.indexOf(key) !== -1 }

    handleJWTChange(jwt) {
        this.props.handleJWTChange(jwt);
    }

    render() {
        return (
            <Fragment>
                <h2>Login</h2>
                <hr></hr>
                <Alert
                    type={this.state.alert.type}
                    message={this.state.alert.message}
                />
                <form className="pt-3" onSubmit={this.handleSubmit}>
                    <InputText
                        title="Email"
                        name="email"
                        type="text"
                        value={this.state.email}
                        handleChange={this.handleChange}
                        className={this.hasErrors("email") ? "is-invalid" : ""}
                        errorDiv={this.hasErrors("email") ? "text-danger" : "d-none"}
                        errorMsg={this.hasErrors("email") ? "Email address required" : "d-none"}
                    />
                    <InputText
                        title="Password"
                        name="password"
                        type="password"
                        value={this.state.password}
                        handleChange={this.handleChange}
                        className={this.hasErrors("password") ? "is-invalid" : ""}
                        errorDiv={this.hasErrors("password") ? "text-danger" : "d-none"}
                        errorMsg={this.hasErrors("password") ? "Password required" : "d-none"}
                    />
                    <button className="btn btn-primary">Login</button>
                </form>
            </Fragment>
        )
    }
}