import React from 'react';
import Loader from "../../loader/Loader";

function Supplies(props) {
    return ( !props.dsLoading && props.districtSupplies !== null? (
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
                {props.districtSupplies.supplies !== null && (
                    props.districtSupplies.supplies.map(supply => (
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
