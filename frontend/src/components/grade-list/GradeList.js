import React, { useState, useEffect } from 'react';
import './grade-list.scss'

function GradeList(props){
    const [gradeList , setgradeList] = useState([]);


    useEffect(() => {
        async function fetchUrl() {
            const schoolID = props.match.params.id;
            //Todo: Check that listID is an integer before parsing it
            const response = await fetch("./api/v1/supply_lists/school/"+schoolID);
            const json = await response.json();
            if(!Array.isArray(gradeList)){
                setgradeList(null);
            }else{
                setgradeList(json);
            }
        }

        fetchUrl();

    },[]);

    return(
        <div className="supply-lists">
            <h2>Select the Supply List you're looking for</h2>
            <ul>
                {gradeList !== null ?
                    (gradeList.map(list => (
                            <li key={list.list_id}>
                                <a href={".#/list/"+list.list_id}>List for grade {list.grade}</a>
                            </li>
                        ))
                    ):(
                        <div>It doesn't look like this school has any supply lists!</div>
                    )
                }
            </ul>
        </div>
    )
}
export default GradeList;