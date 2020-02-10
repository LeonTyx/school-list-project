import React from 'react';
import {NavLink} from 'react-router-dom'
import './dashboard.scss'

function Dashboard() {
    return (
        <div className="dashboard-nav">
            <h3>Dashboard</h3>

            <NavLink activeClassName="active" to={"/dashboard/supplies"}>
                Supplies
            </NavLink>
        </div>
    );

}

export default Dashboard;
