import React, {useEffect, useState} from 'react';
import CardinalToOrdinary from "../cardinal-to-ordinary/CardinalToOrdinary";
import ItemsFoundCounter from "./items-found-counter";

function SupplyList(props) {
    const [supplyList, setSupplyList] = useState(null);

    const [optionalSupplies, setOptionalSupplies] = useState(null);
    const [requiredSupplies, setRequiredSupplies] = useState(null);

    useEffect(() => {
        if(props.list !== null){
            let requiredList = [];
            let optionalList = [];
            const json = props.list;

            json.supply_list.map(supply => !supply.optional ?
                (requiredList.push(supply))
                :
                (optionalList.push(supply))
            );

            setRequiredSupplies(requiredList);
            setOptionalSupplies(optionalList);
            setSupplyList(json);
        }else{
            setSupplyList(null)
        }
    }, [props.list]);

    return (
        supplyList !== null ? (
            <div className="supply-list">
                <ItemsFoundCounter list={requiredSupplies} title={CardinalToOrdinary(props.grade)}/>
                <ItemsFoundCounter list={optionalSupplies} title={"Optional supplies"}/>
            </div>
        ) : (
            <div>Looks like you found an empty supply list</div>
        )
    )
}


export default SupplyList;