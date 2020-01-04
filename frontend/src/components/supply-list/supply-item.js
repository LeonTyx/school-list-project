import React from "react";

class SupplyItem extends React.Component {
    render() {
        return (
            <li className={this.props.optional ? "optional-item" : null}>
                <input
                    type="checkbox"
                    name="supply-item"
                    onChange={this.props.toggleCompletion}
                />
                <h3>{this.props.item}</h3>
                <p>{this.props.desc}</p>
            </li>
        );
    }
}

export default SupplyItem;
