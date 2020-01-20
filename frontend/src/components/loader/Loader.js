import React from 'react';
import './loader.scss'
import LoadingIcon from "../../assets/loader.svg"

function Loader(){
    return(
        <div className="spinner">
            <div className="inner-container">
                <img src={LoadingIcon} alt=""/>
            </div>
        </div>
    )
}
export default Loader;