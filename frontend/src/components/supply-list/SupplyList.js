import React, {useEffect, useState} from 'react';
import CardinalToOrdinary from "../cardinal-to-ordinary/CardinalToOrdinary";
import './supply-list.scss'
import Loader from "../loader/Loader";
import ItemsFoundCounter from "./items-found-counter";

function SupplyList(props) {
    const [loading, setLoading] = useState(true);

    const [supplyList, setSupplyList] = useState(null);

    const [optionalSupplies, setOptionalSupplies] = useState(null);
    const [requiredSupplies, setRequiredSupplies] = useState(null);

    useEffect(() => {
        async function fetchUrl() {
            const listID = props.match.params.id;
            //Todo: Check that listID is an integer before parsing it
            const response = await fetch("./api/v1/supply_lists/" + listID);
            const json = await response.json();
            if (!Array.isArray(json.supply_list)) {
                setSupplyList(null);
                setLoading(false)
            } else {
                let requiredList = [];
                let optionalList = [];

                json.supply_list.map(supply => !supply.optional ?
                    (requiredList.push(supply))
                    :
                    (optionalList.push(supply))
                );

                setRequiredSupplies(requiredList);
                setOptionalSupplies(optionalList);
                setSupplyList(json);

                setLoading(false)
            }
        }

        fetchUrl();
    }, [props.match.params.id]);

    return (
        supplyList !== null ? (
            <div className="supply-list">
                <ItemsFoundCounter list={requiredSupplies} title={CardinalToOrdinary(supplyList.grade)}/>
                <ItemsFoundCounter list={optionalSupplies} title={"Optional supplies"}/>
            </div>
        ) : (
            !loading ? (
                <div>Looks like you found an empty supply list</div>
            ) : (
                <Loader/>
            )
        )
//todo Rewrite to have loading, empty list, and nonexistent list
    )
}


export default SupplyList;