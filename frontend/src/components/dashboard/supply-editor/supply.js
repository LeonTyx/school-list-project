import React, {useEffect, useState} from 'react';
//TODO: Implement better error handling! So far there is no error
//Error handling for non-json responses
function Supply(props) {
    const [isDeleted, setDeletionStatus] = useState(false);

    const [originalText, setOriginalText] = useState({
        "supply_name": props.supply.supply_name,
        "supply_desc": props.supply.supply_desc
    });

    const [currentName, setCurrentName] = useState(props.supply.supply_name);
    const [currentDesc, setCurrentDesc] = useState(props.supply.supply_desc);

    useEffect(() => {
        console.log(originalText)
    }, [originalText]);

    function deleteSupply() {
        sendDeletionRequestion('./api/v1/supplies/supply/'+props.supply.supply_id)
            .then((data) => {
                console.log(data);
                setDeletionStatus(true)
            }
        );

        async function sendDeletionRequestion(url = '') {
            // Default options are marked with *
            const response = await fetch(url, {
                method: 'DELETE',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json'
                },
                redirect: 'follow',
                referrerPolicy: 'no-referrer',
            });

            return await response.json(); // parses JSON response into native JavaScript objects
        }
    }

    function changeSupply() {
        sendChangeRequest('/api/v1/supplies/supply', {
            "supply_id":props.supply.supply_id,
            "supply_name":currentName,
            "supply_desc":currentDesc
        }).then((data) => {
            //Set the original text to be the saved text
            setOriginalText({
                "supply_name": currentName,
                "supply_desc": currentDesc
            });
        });

        async function sendChangeRequest(url = '', data = {}) {
            // Default options are marked with *
            const response = await fetch(url, {
                method: 'PUT',
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

    return (
        !isDeleted &&(
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
                                changeSupply()
                            }}>Save Changes</button>

                            <button onClick={()=>{
                                setCurrentDesc(originalText.supply_desc);
                                setCurrentName(originalText.supply_name);
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
                    <td data-label="Name" className="name">{currentName}</td>
                    <td data-label="Desc" className="desc">{currentDesc}</td>
                </tr>
            )
        )
    )
}

export default Supply;
