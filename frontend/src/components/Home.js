import React, { Component } from 'react';
import Ticket from "./../images/movie-tickets.jpg"
import "./Home.css"

export default class Home extends Component {
    render() {
        return (
            <div className="text-center">
                <h2>This is the Home page</h2>
                <hr />
                <img src={Ticket} alt="movie-ticket" />
                <div className="tickets"></div>
            </div>
        );
    }
}