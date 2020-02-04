import React, {useEffect, useState} from 'react';
import {NavLink} from 'react-router-dom'
import './navbar.scss'

function NavBar() {
    const [schools, setSchools] = useState([null]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchUrl() {
            const response = await fetch("./api/v1/schools");
            const json = await response.json();
            setSchools(json);

            setLoading(false)
        }

        fetchUrl();
    }, []);

    return (
        !loading && (
            schools.schools !== null && schools.district !== "" && (
                <nav id="navigation">
                    {schools.schools.map(school => (
                        <NavLink key={school.school_id} activeClassName="active" to={"/school/" + school.school_id}>
                            {school.name}
                        </NavLink>
                    ))}
                </nav>
            )
        )
    )
}

export default NavBar;