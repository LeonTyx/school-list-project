import React, {useEffect, useState} from 'react';
import Loader from "../components/loader/Loader";
import SupplyList from "../components/supply-list/SupplyList";
import GradeList from "../components/grade-list/GradeList";

import './school-or-list.scss'

function SchoolOrList(props) {
    const [list, setList] = useState(null);
    const [school, setSchool] = useState(null);

    const [listLoading, setListLoading] = useState(false);
    const [gradeLoading, setGradeLoading] = useState(true);

    useEffect(() => {
        setGradeLoading(true);
        fetchSchool();

        async function fetchSchool() {
            const schoolID = props.match.params.schoolID;

            const response = await fetch("/api/v1/supply_lists/school/" + schoolID);
            const json = await response.json();
            if (!Array.isArray(json.supply_lists)) {
                setSchool(null);
            } else {
                setSchool(json);
            }

            setGradeLoading(false);
        }
    }, [props.match.params.schoolID]);

    useEffect(() => {
        if(props.match.params.grade !== undefined){
            setListLoading(true);
            fetchList();
        }

        async function fetchList() {
            const grade = props.match.params.grade;
            const schoolID = props.match.params.schoolID;

            const response = await fetch("/api/v1/supply_lists/school/"+schoolID+"/grade/" + grade);
            const json = await response.json();
            if (!Array.isArray(json.supply_list)) {
                setList(null);
            } else {
                setList(json);
            }

            setListLoading(false);
        }
    }, [props.match.params.grade]);


        return (
        !gradeLoading ? (
            <div className="grade-supply-list">
                <GradeList school={school} schoolID={props.match.params.schoolID}/>
                {props.match.params.grade === undefined && school.school !== null?(
                    <div>Select a grade above </div>
                ):(
                    !listLoading ?
                    (<SupplyList list={list} grade={props.match.params.grade}/>):
                    (<Loader/>)
                )}
            </div>
        ):(
            <Loader/>
        )
    );

}

export default SchoolOrList;
