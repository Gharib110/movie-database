import React, { Component, Fragment } from "react";
import { Link } from "react-router-dom";

export default class OneGenre extends Component {
  state = {
    movies: [],
    isLoaded: false,
    error: null,
    genreName: "",
  };

  componentDidMount() {
    fetch("http://localhost:4000/v1/movies/" + this.props.match.params.id)
      .then((response) => {
        if (response.status !== "200") {
          let err = Error;
          err.message = "Invalid response code: " + response.status;
          this.setState({ error: err });
        }
        return response.json();
      })
      .then((json) => {
        this.setState(
          {
            movies: json.movies,
            isLoaded: true,
            genreName: this.props.location.genreName,
          },
          (error) => {
            this.setState({
              isLoaded: true,
              error,
            });
          }
        );
      });
  }

  render() {
    let { movies, isLoaded, error, genreName } = this.state;

    if (!movies) {
      movies = [];
    }

    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <p>Loading...</p>;
    } else {
      return (
        <Fragment>
          <h2>Genre: {genreName}</h2>

          <div className="list-group">
            {movies.map((m) => (
              <Link
                key={m.id}
                to={`/movies/${m.id}`}
                className="list-group-item list-group-item-action"
              >
                {m.title}
              </Link>
            ))}
          </div>
        </Fragment>
      );
    }
  }
}
