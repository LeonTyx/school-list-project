import React, {useEffect, useState} from 'react';
import './supply-editor.scss'
import Supplies from "./supplies";

function SupplyEditor() {
    const [supplyName, setSupplyName] = useState("");
    const [supplyDesc, setSupplyDesc] = useState("");
    const [districtSupplies, setDS] = useState(null);
    const [dsLoading, setDSLoading] = useState(null);

    useEffect(() => {
        setDSLoading(true);
        fetchDS();

        async function fetchDS() {
            const response = await fetch("/api/v1/supplies/" + 5305400);
            const json = await response.json();
            setDS(json);

            setDSLoading(false);
        }
    }, []);

    function sendData(event) {
        event.preventDefault();

        sendSupply('/api/v1/supplies/supply', {
            "supply_name":supplyName,
            "supply_desc":supplyDesc
        }).then((data) => {
            console.log(data); // JSON data parsed by `response.json()` call

            const newSupplies = districtSupplies;
            newSupplies.supplies.push({
                "supply_id":data,
                "supply_name":supplyName,
                "supply_desc":supplyDesc
            });

            setDS(newSupplies);
            setSupplyName("");
            setSupplyDesc("");
        });

        async function sendSupply(url = '', data = {}) {
            // Default options are marked with *
            const response = await fetch(url, {
                method: 'POST',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json'
                },
                redirect: 'follow',
                referrerPolicy: 'no-referrer',
                body: JSON.stringify(data)
            });
            return await response.json(); // parses JSON response into native JavaScript objects
        }
    }
//../api/v1/supplies/supply
    return (
        <React.Fragment>
            <form onSubmit={sendData}>
                <label>
                    Supply Name:
                    <input
                        type="text"
                        value={supplyName}
                        onChange={e => setSupplyName(e.target.value)}
                    />
                </label>

                <label>
                    Supply Description:
                    <input
                        type="text"
                        value={supplyDesc}
                        onChange={e => setSupplyDesc(e.target.value)}
                    />
                </label>

                <input type="submit" value="Submit" />
            </form>
            <Supplies districtSupplies={districtSupplies}
                      dsLoading={dsLoading}/>
        </React.Fragment>
    );

}

export default SupplyEditor;
