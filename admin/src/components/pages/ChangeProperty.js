import React, { useContext, useEffect, useState } from 'react';
import PropertyForm from '../UI/PropertyForm';

import { GlobalStore } from '../../store/globalStore';


const ChangeProperty = (props) => {
    let { propertyDetails, setProperties, properties } = useContext(GlobalStore);
    let [propertyData, setPropertyData] = useState({});

    useEffect(() => {
        if (propertyDetails) {
            setPropertyData(propertyDetails)
        }
    });

    //submit form
    function sendingData(e) {
        e.preventDefault();
        console.log("Clicked send data button");
        // TODO THINK OF A BETTER WAY OF SHOWING IT'S CHANGING THE PROPERTY
        props.sendDataToBackend(e.target, true);
    }

    return (<div>
        <PropertyForm
            sendingData={sendingData}
            propertyDetails={propertyData}
            propertyEdit={true}
        />
    </div>);
}


export default ChangeProperty;