import React, { useState, useEffect } from 'react';
import './schools.scss'
import LandingHero from '../../assets/back-to-school.png'

function Schools(){
    const [schools , setSchools] = useState([null]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchUrl() {
            const response = await fetch("./schools");
            const json = await response.json();
            setSchools(json);
            setLoading(false)
        }

        fetchUrl();
    }, []);
    //todo: add "no schools" in case of empty school list
    return(
        !loading ?(
            <div className="schools">
                <img src={LandingHero} alt="Child ready for first day of school"/>
                <h2>Schools in {schools.district}</h2>

                <ul>
                    {schools.schools.map(school => (
                        <li key={school.school_id}>
                            <a href={".#/lists/"+school.school_id}>List for {school.name}</a>
                        </li>
                    ))}
                </ul>
            </div>
        ):(
            <div>Loading...</div>
        )
    )
}
export default Schools;