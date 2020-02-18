import React, {useEffect, useState} from 'react';

function Supply(props) {
    const [originalText, setOriginalText] = useState({
        "supply_name": props.supply.supply_name,
        "supply_desc": props.supply.supply_desc
    });

    const [currentName, setCurrentName] = useState(props.supply.supply_name);
    const [currentDesc, setCurrentDesc] = useState(props.supply.supply_desc);

    useEffect(() => {
        console.log(originalText)
    }, [originalText]);

    function deleteSupply(event) {
        console.log(props.supply.supply_id, "deleted")
    }

    return (
        props.editingMode ? (
            <tr>
                <td data-label="ID" className="ID">{props.supply.supply_id}</td>
                <td data-label="Name" className="name">
                    <input value={currentName}
                           onChange={e => setCurrentName(e.target.value)}/>
                </td>
                <td data-label="Desc" className="desc">
                    <input value={currentDesc}
                           onChange={e => setCurrentDesc(e.target.value)}/>
                </td>
                {((currentName !== originalText.supply_name) ||
                    (currentDesc !== originalText.supply_desc)) ? (
                    <td data-label="Edit">
                        <button onClick={()=>{
                            setOriginalText({
                                "supply_name": currentName,
                                "supply_desc": currentDesc
                            })
                        }}>Save Changes</button>

                        <button onClick={()=>{
                            setCurrentDesc(currentDesc);
                            setCurrentName(currentName);
                        }}>Cancel</button>
                    </td>
                    ):(
                        <td data-label="Edit">
                            <button onClick={deleteSupply}>Delete</button>
                        </td>
                    )
                }
            </tr>
        ):(
            <tr>
                <td data-label="ID" className="ID">{props.supply.supply_id}</td>
                <td data-label="Name" className="name">{props.supply.supply_name}</td>
                <td data-label="Desc" className="desc">{props.supply.supply_desc}</td>
            </tr>
        )
    );

}

export default Supply;
