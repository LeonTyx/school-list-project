import React, {useEffect, useState} from 'react';
import {NavLink} from 'react-router-dom'
import CardinalToOrdinary from "../cardinal-to-ordinary/CardinalToOrdinary";

function GradeList(props) {
    const [grades, setGrades] = useState(null);
    const [schoolName, setSchoolName] = useState(null);
    const [educationStage, setEducationStage] = useState(null);

    useEffect(() => {
        if(props.school!== null){
            setGrades(props.school.supply_lists);
            setSchoolName(props.school.school_name);
            setEducationStage(props.school.education_stage);
        }
    }, [props.school]);

    return(
        props.school !== null ?(
            <div className="school-info">
                <div className="school-header">
                    <h2>{schoolName}</h2>
                    <h3>{educationStage}</h3>
                </div>

                <p>Select grade</p>
                <div className="list-nav">
                    <ul className="list-picker">
                        {grades !== null && (
                            grades.map(grade => (
                                    <li>
                                        <NavLink key={grade.grade} activeClassName="active" to={"/school/"+ props.schoolID +"/" + grade.grade}>
                                            {CardinalToOrdinary(grade.grade)}
                                        </NavLink>
                                    </li>
                                )
                            ))}
                    </ul>
                </div>
            </div>
        ):(
            <div> You've stumbled upon a school without any supply lists! </div>
        )
    )
}

export default GradeList;