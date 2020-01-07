import React, { useState, useEffect } from 'react';
import CardinalToOrdinary from "../cardinal-to-ordinary/CardinalToOrdinary";
import './supply-list.scss'
import SupplyItem from "./supply-item";

function SupplyList(props){
    const [loading, setLoading] = useState(true);

    const [supplyList , setSupplyList] = useState(null);

    const [requiredItemsCompleted , setRequiredItemsCompleted] = useState(0);
    const [optionalItemsCompleted , setOptionalItemsCompleted] = useState(0);

    const [optionalSupplies , setOptionalSupplies] = useState(null);
    const [requiredSupplies , setRequiredSupplies] = useState(null);


    function handleRequiredCompletion(event) {
        let checked = event.target.checked;
        let completed = requiredItemsCompleted;

        if (checked) {
            completed++;
        } else {
            completed--;
        }
        setRequiredItemsCompleted(completed)
    }

    function handleOptionalCompletion(event) {
        let checked = event.target.checked;
        let completed = optionalItemsCompleted;

        if (checked) {
            completed++;
        } else {
            completed--;
        }
        setOptionalItemsCompleted(completed)
    }

    useEffect(() => {
        async function fetchUrl() {
            const listID = props.match.params.id;
            //Todo: Check that listID is an integer before parsing it
            const response = await fetch("./api/v1/supply_lists/"+listID);
            const json = await response.json();
            if(!Array.isArray(json.supply_list)){
                setSupplyList(null);
                setLoading(false)
            }else{
                let requiredList = [];
                let optionalList =[];

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
    },[props.match.params.id]);

    return(
        supplyList !== null ? (
            <div className="supply-list">
                <div className="title">
                    <h2>{CardinalToOrdinary(supplyList.grade)}</h2>
                    <div className="completed-items">
                        {requiredItemsCompleted} of {requiredSupplies.length} items found
                    </div>
                </div>

                <ul>
                    {requiredSupplies.map(item => (
                        <SupplyItem
                            key={item.supply_id}
                            item={item.name}
                            desc={item.desc}
                            optional={item.optional}
                            toggleCompletion={handleRequiredCompletion}
                        />
                    ))}
                </ul>


                <div className="title">
                    <h4>Optional supplies</h4>
                    <div className="completed-items">
                        {optionalItemsCompleted} of {optionalSupplies.length} items found
                    </div>
                </div>

                <ul className="optional">
                    {optionalSupplies.map(item => (
                        <SupplyItem
                            key={item.supply_id}
                            item={item.name}
                            desc={item.desc}
                            optional={item.optional}
                            toggleCompletion={handleOptionalCompletion}
                        />
                    ))}
                </ul>
            </div>
        ):(
            !loading ? (
                <div>Looks like you found an empty supply list</div>
            ):(
                <div className={"loading-text"}>Loading...</div>
            )
        )
//todo Rewrite to have loading, empty list, and nonexistent list
    )
}


export default SupplyList;