import React, {useEffect, useState} from 'react';
import './schools.scss'
import LandingHero from '../../assets/back-to-school.svg'
import Loader from "../loader/Loader";

function Schools() {
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
        !loading ? (
            schools.schools !== null && schools.district !== "" ? (
                <div className="schools">
                    <img src={LandingHero} alt="Child ready for first day of school"/>
                    <h2>Schools in {schools.district}</h2>

                    <ul>
                        {schools.schools.map(school => (
                            <li key={school.school_id}>
                                <a href={".#/lists/" + school.school_id}>List for {school.name}</a>
                            </li>
                        ))}
                    </ul>
                </div>
            ) : (
                <div>Looks like this district doesn't have any schools</div>
            )

        ) : (
            <Loader/>
        )
    )
}

export default Schools;