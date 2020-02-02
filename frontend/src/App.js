import React, {useState} from 'react';
import './app.scss'

import Header from "./components/header/Header";
import Schools from "./components/schools/Schools";
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import GradeList from "./components/grade-list/GradeList";
import SupplyList from "./components/supply-list/SupplyList";
import NavBar from "./components/navbar/NavBar";

//Todo: rewrite entire frontend session logic. This terrifies me.
function App() {
    const [storageStatus, setStorageStatus] = useState(null);
    const [user, setUser] = useState({
        name: null,
        email: null,
        isLoggedIn: false
    });

    return (
        <Router>
            <Header user={user}
                    setUser={setUser}
                    storageStatus={storageStatus}
                    setStorageStatus={setStorageStatus}/>
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

export default App;
