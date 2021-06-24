import React, { Component, Fragment } from 'react'

export default class ShowOneMovie extends Component{
    state = { movie : {}};

    componentDidMount() {
        this.setState({
            movie: {
                id: this.props.match.params.id,
                title: "The Shawshank Redemption", 
                runtime: 160
            }
        })
    }

    render() {
        return(
            <Fragment>
                <h2>Id: {this.state.movie.id}</h2>
                <h2></h2>
                <h2></h2>
                <table className="table table-compact table-light">
                    <thead></thead>
                    <tbody>
                        <tr>
                            <td><strong>Title:</strong></td>
                            <td><strong>{this.state.movie.title}</strong></td>
                        </tr>
                        <tr>
                            <td><strong>Runtime:</strong></td>
                            <td><strong>{this.state.movie.runtime} min</strong></td>
                        </tr>
                    </tbody>
                </table>
            </Fragment>
        )
    }
}