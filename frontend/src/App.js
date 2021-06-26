import React, { Component, Fragment } from 'react';
import { HashRouter as Router, Switch, Route, Link } from 'react-router-dom'

import Admin from './components/Admin';
import Genres from './components/Genres';
import Home from './components/Home';
import Movies from './components/Movies';
import ShowOneMovie from './components/ShowOneMovie';
import MoviesByGenre from './components/MoviesByGenre';
import EditMovie from './components/EditMovie';
import Login from './components/Login';





export default class App extends Component {

  constructor(props) {
    super(props);
    this.state = {
      jwt: "",
    }
    this.handleJWTChange = this.handleJWTChange.bind(this);

  }

  handleJWTChange = jwt => {
    this.setState({ jwt: jwt });
  }

  logout = () => {
    this.setState({ jwt: "" })
  }

  render() {
    return (
      <Router>
        <div className="container">
          <div className="row">
            <div className="col mt-3">
              <h1 className="mt-3">
                Go Watch a Movie
            </h1>
            </div>
            <div className="col mt-3 text-end">
              {this.state.jwt === "" ?
                <Link to="/login">Login</Link> :
                <Link to="/logout" onClick={this.logout}>Logout</Link>
              }
            </div>
            <hr className="mb-3"></hr>


          </div>

          <div className="row">
            <div className="col-md-2">
              <nav>
                <ul className="list-group">
                  <li className="list-group-item">
                    {/* <a href="/">Home</a> */}
                    <Link to="/">Home</Link>

                  </li>
                  <li className="list-group-item">
                    <Link to="/movies">Movies</Link>
                  </li>

                  <li className="list-group-item">
                    <Link to="/genres">Genres</Link>
                  </li>
                  {this.state.jwt !== "" && (
                    <Fragment>
                      <li className="list-group-item">
                        <Link to="/admin/movies/0">Add movie</Link>
                      </li>

                      <li className="list-group-item">
                        <Link to="/admin">Manage Cataloge</Link>
                      </li>
                    </Fragment>
                  )}
                </ul>
                <pre>
                  {JSON.stringify(this.state, null, 3)}
                </pre>
              </nav>
            </div>

            <div className="col-md-10">
              <Switch>
                <Route path="/login" component={props => <Login {...props} handleJWTChange={this.handleJWTChange}/>}/>

                <Route path="/movies/:id" component={ShowOneMovie}>
                </Route>

                <Route path="/movies">
                  <Movies />
                </Route>


                <Route path="/genres/:id" component={MoviesByGenre}>
                </Route>
                <Route exact path="/genres">
                  <Genres />
                </Route>

                <Route path="/admin/movies/:id" component={EditMovie}>
                </Route>

                <Route path="/admin">
                  <Admin />
                </Route>

                <Route exact path="/">
                  <Home />
                </Route>

              </Switch>
            </div>
          </div>
        </div>
      </Router>

    );
  }
}
