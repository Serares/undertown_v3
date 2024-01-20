import React from 'react';


import PropertyForm from '../UI/PropertyForm';


function AddProperty(props) {

    //submit form
    function sendingData(e) {
        e.preventDefault();
        console.log("Clicked send data button");
        // get all the data and add it to fields
        props.sendDataToBackend(e.target, false);
    }

    return (
        <div className="addProperty">
            {<PropertyForm
                sendingData={sendingData}
                propertyEdit={false}
            />}
        </div>
    )

}

export default AddProperty;