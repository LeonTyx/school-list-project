import React, {useState} from 'react';
import './app.scss'

import Header from "./components/header/Header";
import Schools from "./components/schools/Schools";
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import NavBar from "./components/navbar/NavBar";
import SchoolOrList from "./school-or-list/school-or-list";
import NotFound from "./components/not-found/not-found";
import {ProtectedRoute} from "./components/protected-routes/ProtectedRoute";
import Dashboard from "./components/dashboard/dashboard";
import SupplyEditor from "./components/dashboard/supply-editor/supply-editor";

function App() {
    const [storageStatus, setStorageStatus] = useState(null);
    const [user, setUser] = useState({
        name: null,
        email: null,
        isLoggedIn: false,
    });

    return (
        <Router>
            <Header user={user}
                    setUser={setUser}
                    storageStatus={storageStatus}
                    setStorageStatus={setStorageStatus}/>
            <main>
                <Switch>
                    <ProtectedRoute exact path="/dashboard" component={Dashboard} user={user} />
                    <ProtectedRoute exact path="/dashboard/supplies" component={SupplyEditor} user={user} />

                    <Route exact path="/" component={Schools}/>

                    <Route path="/school/:schoolID/:grade?" component={SchoolOrList}/>
                    <Route component={NotFound} />
                </Switch>
            </main>
            <NavBar isLoggedIn={user.isLoggedIn}/>
        </Router>
    );

}

export default App;
