import React, { useState, useEffect } from 'react';
import './grade-list.scss'
import Loader from "../loader/Loader";

function GradeList(props){
    const [gradeList , setGradeList] = useState([]);
    const [loading, setLoading] = useState(true);


    useEffect(() => {
        setLoading(true);

        async function fetchUrl() {
            const schoolID = props.match.params.id;
            //Todo: Check that listID is an integer before parsing it
            const response = await fetch("./api/v1/supply_lists/school/"+schoolID);
            const json = await response.json();

            setGradeList(json);

            setLoading(false)
        }

        fetchUrl();

    },[props.match.params.id]);

    return !loading ? (
        gradeList !== null ? (
            <div className="school">
                <div className="grade-list-header">
                    <h2>
                        {gradeList.school_name !== null
                            ? gradeList.school_name
                            : "This school has no name"}
                    </h2>
                    <div className="education-stage">{gradeList.education_stage}</div>
                </div>
                <ul className="grade-list">
                    {gradeList.supply_lists !== null ? (
                        gradeList.supply_lists.map(list => (
                            <li key={list.list_id}>
                                <a href={".#/list/" + list.list_id}>
                                    List for grade {list.grade}
                                </a>

                                <div className="school-year">
                                    {list.starting_year} to {list.ending_year}
                                </div>
                            </li>
                        ))
                    ) : (
                        <li>It doesn't look like this school has any supply lists!</li>
                    )}
                </ul>
            </div>
        ) : (
            <div>Doesn't look like this school exists</div>
        )
    ) : (
        <Loader/>
    );
}
export default GradeList;