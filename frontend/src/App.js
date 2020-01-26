import React from 'react';
import './app.scss'

import Header from "./components/header/Header";
import Schools from "./components/schools/Schools";
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import GradeList from "./components/grade-list/GradeList";
import SupplyList from "./components/supply-list/SupplyList";
import NavBar from "./components/navbar/NavBar";

//Todo: rewrite entire frontend session logic. This terrifies me.
class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isLoggedIn: false,
            name: null,
            email: null,
            storageStatus: null
        };
    }

    componentDidMount() {
        console.log("Component mounted");

        //Every 30 seconds, check if session needs to be refresh
        setInterval(this.refreshSession, 30000);

        if (this.cookieExists("session")) {
            console.log("Cookie exists: Session");
            this.localStorageUpdated();
        }
        //Start by checking if a user session cookie exists
        this.checkSession();

        //User localstorage to save user data across tabs!
        //Uses event listener to check if data has changed and updates
        // states in all tabs. Hopefully.
        if (typeof window !== 'undefined') {
            this.setState({storageStatus: !!localStorage.getItem('session')});

            window.addEventListener('storage', this.localStorageUpdated)
        }
    }

    refreshSession = () => {
        if (this.cookieExists("session")) {
            if (localStorage.getItem("isRefreshing") === "false" ||
                localStorage.getItem("isRefreshing") === null) {
                const lastLogin = new Date(localStorage.getItem("lastRefresh"));
                const currentDate = new Date();

                //If there has been more than 25 minutes since login, refresh session
                if ((currentDate.getTime() - lastLogin.getTime()) / 1000 > 1500) {
                    localStorage.setItem("isRefreshing", true);

                    fetch("./oauth/v1/refresh")
                        .then(respone => {
                            console.log("Session refreshed");
                            localStorage.setItem("lastRefresh", new Date());
                            localStorage.setItem("isRefreshing", false);
                        })
                }
            } else {
                console.log("Cannot initiate session refresh. A session refresh was already started");
            }
        } else {
            this.setState({
                email: null,
                name: null,
                isLoggedIn: false
            });

            localStorage.removeItem("name");
            localStorage.removeItem("email");
        }
    };

    checkSession = () => {
        console.log("attempting to get user data");
        if (this.cookieExists("session")) {
            fetch('./oauth/v1/profile')
                .then(response => {
                    if (response.status === 414) {
                        localStorage.removeItem('name');
                        localStorage.removeItem('email');
                        this.setState({
                            name: null,
                            email: null,
                            isLoggedIn: false
                        })
                    }
                    return response.json();
                })
                .then(profile => {
                    localStorage.setItem('name', profile.name);
                    localStorage.setItem('email', profile.email);
                    this.setState({
                        name: profile.name,
                        email: profile.email,
                        isLoggedIn: true
                    })
                });
        }
    };

    localStorageUpdated = () => {
        console.log("Local storage updated");
        const checkSession = this.checkSession;
        //Get user data from local storage
        const name = localStorage.getItem('name');
        const email = localStorage.getItem('email');

        console.log(name, email);
        if (name === null && email === null) {
            console.log("User not logged in");
            //There is no user logged in
            this.setState({name: null, email: null, isLoggedIn: false})
        } else if (name !== null && email !== null) {
            console.log("User logged in");

            //There is a user logged in
            this.setState({name: name, email: email, isLoggedIn: true})
        } else {
            console.log("User information in localstorage have become decoupled. Attempting to contact server about current user session");
            checkSession()
        }
        //todo: Use switch instead of this boolean nightmare
    };

    logout = () => {
        console.log("Logging out");
        fetch('./oauth/v1/logout')
            .then(response => {
                localStorage.removeItem('name');
                localStorage.removeItem('email');
                this.setState({
                    name: null,
                    email: null,
                    isLoggedIn: false
                })
            })
    };

    setLoginTime = () => {
        localStorage.setItem("lastRefresh", new Date())
    };

    cookieExists = (name) => {
        var dc = document.cookie;
        var prefix = name + "=";
        var begin = dc.indexOf("; " + prefix);
        if (begin === -1) {
            begin = dc.indexOf(prefix);
            if (begin !== 0) return null;
        } else {
            begin += 2;
            var end = document.cookie.indexOf(";", begin);
            if (end === -1) {
                end = dc.length;
            }
        }
        // because unescape has been deprecated, replaced with decodeURI
        //return unescape(dc.substring(begin + prefix.length, end));
        return decodeURI(dc.substring(begin + prefix.length, end)) != null;
    };

    render() {
        return (
            <Router>
                <Header isLoggedIn={this.state.isLoggedIn} name={this.state.name} logout={this.logout}
                        setLoginTime={this.setLoginTime}/>
                <main>
                    <Switch>
                        <Route exact path="/" component={Schools}/>

                        <Route path="/lists/:id" component={GradeList}/>
                        <Route path="/list/:id" component={SupplyList}/>

                    </Switch>
                </main>

                <NavBar/>
            </Router>
        );
    }

}

export default App;
