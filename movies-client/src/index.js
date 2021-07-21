import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import AppHeader from './AppHeader';
import AppFooter from './AppFooter';
import AppContent from './AppContent';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import './index.css';

class App extends Component {

  constructor(props) {
    super(props);
    this.handlePostChange = this.handlePostChange.bind(this);
    this.state = {posts: []};
  }

  handlePostChange(posts) {
    this.setState({posts: posts});
  }

  render() {
    const myProps = {
      title: "My Cool App!",
      subject: "My subject",
      favourite_color: "red",
    }

    return (
      <div className="app">
        <AppHeader {...myProps} posts={this.state.posts} handlePostChange={this.handlePostChange} />
        <AppContent handlePostChange={this.handlePostChange} posts={this.state.posts} />
        <AppFooter />
      </div>
    );
  }
}


ReactDOM.render(<App />, document.getElementById('root'));
