import React from 'react';
import { HashRouter as Router, Switch, Route, Link } from 'react-router-dom'

import Admin from './components/Admin';
import Genres from './components/Genres';
import Home from './components/Home';
import Movies from './components/Movies';
import ShowOneMovie from './components/ShowOneMovie';
import MoviesByGenre from './components/MoviesByGenre';
import EditMovie from './components/EditMovie';




export default function App() {
  return (
    <Router>
      <div className="container">
        <div className="row">
          <h1 className="mt-3">
            Go Watch a Movie
        </h1>
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

                <li className="list-group-item">
                  <Link to="/admin/movies/0">Add movie</Link>
                </li>

                <li className="list-group-item">
                  <Link to="/admin">Manage Cataloge</Link>
                </li>
              </ul>
            </nav>
          </div>

          <div className="col-md-10">
            <Switch>


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
