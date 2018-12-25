import React, { Component } from 'react';
import logo from './twitter_logo.svg';
import './bootstrap.min.css';
import './App.css';

import 'whatwg-fetch';

class App extends Component {
    // TODO websocket URL should get from this.props
    constructor(props) {
	super(props);
	this.state = {
	    value: '',
	    list:[],
	    ws: new WebSocket("ws://127.0.0.1:3000/api/filter")
	};
    };

    logging(logline) {
	console.log(logline);
    };

    
    setupWebsocket() {
        let websocket = this.state.ws;

        websocket.onopen = () => {
          this.logging('Websocket connected');
        };

        websocket.onmessage = (evt) => {
	    var obj = JSON.parse(evt.data);
	    this.logging(obj);
	    
	    this.setState((state) => ({ list: state.list.concat(obj)}));
        };

        websocket.onclose = () => {
          this.logging('Websocket disconnected');
        }
    };

    componentDidMount() {
      this.setupWebsocket();
    };

    componentWillUnmount() {
      let websocket = this.state.ws;
      websocket.close();
    };

    renderItems() {
	const items = [];
	this.state.list.forEach(item => items.push(<li>{item.text} - <b>{item.count}</b></li>));
	return items
    }
  
    handleClickEvent = (e) =>{
	e.preventDefault();
	let websocket = this.state.ws;
	
	websocket.send(this.state.value);
	
   };
    
    handleChange = (event) => {
	this.setState({value: event.target.value});
    };
    
  render() {
      return (
	  <div>
	      <div className="container-fluid">
	        <nav className="navbar navbar-default">
	          <div className="container">
	            <div className="navbar-header">
	              <div className="vertical-center">
	                <img src={logo} className="pull-left banner img-rounded" width="36" alt="logo" />
                        <a className="navbar-brand" href="javascript:void(0)">Twitter Stream API</a>
	              </div>
	            </div>
	          </div>
	        </nav>
	 
	      <div className="container-fluid">
  	       
	  <form className="form-horizontal">
	     <fieldset>
                 <div className="form-group">
	            <div className="col-lg-11">
          <input type="text" autoComplete="off" className="form-control" id="inputKeyword" placeholder="Enter a keyword to track"
	   value={this.state.value}
           onChange={this.handleChange.bind(this)}
	  />
	            </div>
	            <a type="submit" className="btn btn-primary col-lg-1" onClick={this.handleClickEvent.bind(this)}>Query</a>
	         </div>

	     </fieldset>
	      </form>
	      
	    <ul className="list">
              {this.renderItems()}
            </ul>
	  
	    </div>
	  </div>
	</div>
      ); 
  }
}

export default App;
