import React, {useEffect} from 'react'
import './header.scss'

function Header(props) {
    useEffect(() => {
        console.log("Component mounted");
        //Every 30 seconds, check if session needs to be refresh
        setInterval(refreshSession, 30000);

        if (cookieExists("session")) {
            console.log("Cookie exists: Session");
            localStorageUpdated();
        }
        //Start by checking if a user session cookie exists
        checkSession();

        //User localstorage to save user data across tabs!
        //Uses event listener to check if data has changed and updates
        // states in all tabs. Hopefully.
        if (typeof window !== 'undefined') {
            props.setStorageStatus(!!localStorage.getItem('session'));

            window.addEventListener('storage', localStorageUpdated)
        }


        function refreshSession(){
            if (cookieExists("session")) {
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
                props.setUser({
                    name: null,
                    email: null,
                    isLoggedIn: false
                });


                localStorage.removeItem("name");
                localStorage.removeItem("email");
            }
        }

        function checkSession(){
            console.log("attempting to get user data");
            if (cookieExists("session")) {
                fetch('./oauth/v1/profile')
                    .then(response => {
                        if (response.status === 414) {
                            localStorage.removeItem('name');
                            localStorage.removeItem('email');
                            props.setUser({
                                name: null,
                                email: null,
                                isLoggedIn: false
                            });
                        }
                        return response.json();
                    })
                    .then(profile => {
                        localStorage.setItem('name', profile.name);
                        localStorage.setItem('email', profile.email);

                        props.setUser({
                            name: profile.name,
                            email: profile.email,
                            isLoggedIn: true
                        });
                    });
            }
        }

        function localStorageUpdated() {
            console.log("Local storage updated");
            //Get user data from local storage
            const name = localStorage.getItem('name');
            const email = localStorage.getItem('email');

            console.log(name, email);
            if (name === null && email === null) {
                console.log("User not logged in");
                //There is no user logged in
                props.setUser({
                    name: null,
                    email: null,
                    isLoggedIn: false
                });

            } else if (name !== null && email !== null) {
                console.log("User logged in");

                //There is a user logged in
                props.setUser({
                    name: name,
                    email: email,
                    isLoggedIn: true
                });
            } else {
                console.log("User information in localstorage have become decoupled. Attempting to contact server about current user session");
                checkSession()
            }
            //todo: Use switch instead of this boolean nightmare
        }

    },[]);

    function logout() {
        console.log("Logging out");
        fetch('./oauth/v1/logout')
            .then(response => {
                localStorage.removeItem('name');
                localStorage.removeItem('email');
                props.setUser({
                    name: null,
                    email: null,
                    isLoggedIn: false
                });
            })
    }

    function setLoginTime(){
        localStorage.setItem("lastRefresh", new Date())
    }

    function cookieExists(name) {
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
    }

    return (
        <header>
            <a href="./#" className="logo">Back to School Simplified</a>
            {props.user.isLoggedIn ? (
                <div>Hello, {props.user.name}!. <button onClick={logout} className="logout-button">Logout here</button>
                </div>
            ) : (
                <div><a href={"./oauth/v1/login"} onClick={setLoginTime}>Editor Login</a></div>
            )}
        </header>
    )
}

export default Header;