import React, {useEffect, useState} from 'react';
import SupplyItem from "./supply-item";

function ItemsFoundCounter(props) {
    const [suppliesFound, setSuppliesFound] = useState(0);
    const [supplies, setSupplies] = useState(null);

    useEffect(() => {
        setSupplies(props.list)
    },[props]);

    function handleCompletion(event) {
        let checked = event.target.checked;
        let completed = suppliesFound;

        if (checked) {
            completed++;
        } else {
            completed--;
        }
        setSuppliesFound(completed)
    }

    return (supplies !== null && (
        <div>
            <div className="supply-list-header">
                <div className="title">{props.title}</div>
                <div className="completed-items">
                    {suppliesFound} of {supplies.length} items found
                </div>
            </div>

            <ul>
                {supplies.map(supply => (
                    <SupplyItem
                        key={supply.supply_id}
                        item={supply.name}
                        desc={supply.desc}
                        optional={supply.optional}
                        amount={supply.amount}
                        toggleCompletion={handleCompletion}
                    />
                ))}
            </ul>
        </div>
    ))
}


export default ItemsFoundCounter;