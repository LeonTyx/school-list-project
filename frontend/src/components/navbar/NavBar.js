import React from 'react'
import './nav.scss'

function NavBar(props){
    return(
        <header>
            <a href="./#" className="logo">Back to School Simplified</a>
            {props.isLoggedIn ? (
                <div>Hello, {props.name}!. <button onClick={props.logout} className="logout-button">Logout here</button></div>
            ):(
                <div><a href={"./oauth/v1/login"} onClick={props.setLoginTime}>Editor Login</a></div>
            )}
        </header>
    )
}
export default NavBar;