import React from 'react';
import Loader from "../../loader/Loader";
import Supply from "./supply";

function Supplies(props) {
    return ( !props.dsLoading && props.districtSupplies !== null? (
        <div className="supplies">
            <table>
                <caption>Supplies for {props.districtID}</caption>
                <thead>
                    <tr>
                        <th scope="col" className="ID">ID</th>
                        <th scope="col" className="name">Name</th>
                        <th scope="col" className="desc">Description</th>
                    </tr>
                </thead>

                <tbody>
                {props.districtSupplies.supplies !== null && (
                    props.districtSupplies.supplies.map(supply => (
                        <Supply supply={supply} editingMode={props.editingMode}/>
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
