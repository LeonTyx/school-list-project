import React from "react";
import { Route } from "react-router-dom";
import LoginRequired from "../login-required/LoginRequired";

export const ProtectedRoute = ({ component: Component,user, ...rest }) => {
    return (
        <Route
            {...rest}
            render={props => {
                if (user.isLoggedIn) {
                    return <Component {...props} />;
                } else {
                    return <LoginRequired/>
                }
            }}
        />
    );
};
