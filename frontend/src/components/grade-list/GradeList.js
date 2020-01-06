import React, { useState, useEffect } from 'react';
import './grade-list.scss'

function GradeList(props){
    const [gradeList , setGradeList] = useState([]);
    const [loading, setLoading] = useState(true);


    useEffect(() => {
        async function fetchUrl() {
            const schoolID = props.match.params.id;
            //Todo: Check that listID is an integer before parsing it
            const response = await fetch("./api/v1/supply_lists/school/"+schoolID);
            const json = await response.json();

            setGradeList(json);

            setLoading(false)
        }

        fetchUrl();

    },[]);

    return(
        !loading ?(
            gradeList !== null ? (
                <div className="supply-lists">
                    <h2>
                        {gradeList.school_name !== null ?
                            (gradeList.school_name)
                            :
                            ("This school has no name")
                        }
                    </h2>
                    <h3>Select the Supply List you're looking for</h3>
                    <ul>
                        {gradeList.supply_lists !== null ?
                            (gradeList.supply_lists.map(list => (
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
            ):(
                <div>Doesn't look like this school exists</div>
            )
        ):(
            <div>Loading</div>
        )

    )
}
export default GradeList;