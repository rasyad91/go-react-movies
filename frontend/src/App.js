import React from 'react';
import { HashRouter as Router, Switch, Route, Link, useParams, useRouteMatch } from 'react-router-dom'

import Admin from './components/Admin';
import Home from './components/Home';
import Movies from './components/Movies';
import Categories from './components/Categories';
import ShowOneMovie from './components/ShowOneMovie';



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
                  <Link to="/categories">By Category</Link>
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

              <Route exact path="/categories">
                <Category />
              </Route>

              <Route exact path="/categories/drama" render={(props) => <Categories {...props} title={"Drama"} />}/>
              <Route exact path="/categories/comedy" render={(props) => <Categories {...props} title={"Comedy"} />}/>

              

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

function Category() {
  let { path, url } = useRouteMatch();

  return(
      <div>
        <h2>Categories</h2>
        <ul>
          <li>
            <Link to={`${path}/drama`}>Drama path</Link>
          </li>
          <li>
            <Link to={`${url}/drama`}>Drama url</Link>
          </li>
          <li>
            <Link to={`${path}/comedy`}>Comedy</Link>
          </li>
          <li>
            <Link to={`${path}/horror`}>Horror</Link>
          </li>
        </ul>
      </div>
  )
}