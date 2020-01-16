import React from "react";

class SupplyItem extends React.Component {
    render() {
        return (
            <li className={this.props.optional ? "optional-item" : null}>
                <label className="supply-item"><h3>{this.props.item}</h3>
                    <input type="checkbox"
                           name="supply-item"
                           onChange={this.props.toggleCompletion}/>
                    <span className="checkmark"/>
                    <p>{this.props.desc}</p>
                </label>
            </li>
        );
    }
}

export default SupplyItem;
