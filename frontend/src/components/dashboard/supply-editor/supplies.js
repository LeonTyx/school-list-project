import React, {useEffect, useState} from 'react';
import Loader from "../../loader/Loader";

function Supplies(props) {
    const [districtSupplies, setDS] = useState(null);
    const [dsLoading, setDSLoading] = useState(null);

    useEffect(() => {
        setDSLoading(true);
        fetchDS();

        async function fetchDS() {
            const response = await fetch("/api/v1/supplies/" + props.districtID);
            const json = await response.json();
            setDS(json);

            setDSLoading(false);
        }
    }, [props.districtID]);

    return ( !dsLoading && districtSupplies !== null? (
        <div className="supplies">
            <table>
                <caption>Supplies for {props.districtID}</caption>
                <thead>
                    <tr>
                        <th scope="col">ID</th>
                        <th scope="col">Name</th>
                        <th scope="col">Description</th>
                    </tr>
                </thead>

                <tbody>
                {districtSupplies.supplies !== null && (
                    districtSupplies.supplies.map(supply => (
                        <tr>
                            <td data-label="ID">{supply.supply_id}</td>
                            <td data-label="Name">{supply.supply_name}</td>
                            <td data-label="Desc">{supply.supply_desc}</td>
                        </tr>
                        ))
                    )
                }
                </tbody>
            </table>
        </div>
    ):(
        <Loader/>
    ));
}

export default Supplies;
