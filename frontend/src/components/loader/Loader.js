import React from 'react';
import './loader.scss'
import LoadingIcon from "../../assets/loader.svg"

function Loader(){
    return(
        <div className="spinner">
            <img src={LoadingIcon} alt=""/>
        </div>
    )
}
export default Loader;