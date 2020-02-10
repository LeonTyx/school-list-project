import React, {useState} from 'react';
import './supply-editor.scss'

function SupplyEditor() {
    const [supplyName, setSupplyName] = useState("");
    const [supplyDesc, setSupplyDesc] = useState("");

    function sendData(event) {
        event.preventDefault();
        console.log("Supply Name:", supplyName);
        console.log("Supply Desc:", supplyDesc)
    }
//../api/v1/supplies/supply
    return (
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
    );

}

export default SupplyEditor;
