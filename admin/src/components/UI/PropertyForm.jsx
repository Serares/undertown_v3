import React, { useState, useEffect, useContext } from 'react';
import { Button, Form } from 'react-bootstrap';
import Input from './Input';
import classes from './PropertyForm.module.css';
import { GlobalStore } from '../../store/globalStore';

const PropertyForm = (props) => {

    let { setIsPropertyEdit, adminAccounts, propertyFields } = useContext(GlobalStore);
    let [formData, setFormData] = useState(propertyFields);
    let { propertyDetails, propertyEdit } = props;

    useEffect(() => {
        setIsPropertyEdit(propertyEdit);
    })

    useEffect(() => {
        if (adminAccounts) {
            generateContactPersons(adminAccounts);
        }
        // propertyDetails is defined only from ChangeProperty
        if (propertyDetails) {
            getPropertyDataValuesToFormFields(propertyDetails);
        }
    }, [propertyDetails, adminAccounts])

    useEffect(() => {
        window.postMessage({
            scope: "PROPERTY_FORM",
            data: "ELEMENT_MAP",
            loaded: true
        });
    }, []);

    /**
     * @param {accounts}: User[]
     */
    function generateContactPersons(accounts) {
        // {display: "firstName + lastName", value: userID}
        let newFormData = { ...formData };
        let arrayOfContacts = [];
        arrayOfContacts = accounts.map((acc, index) => {
            return {
                display: `${acc["firstName"]} ${acc["lastName"]}`,
                value: acc["_id"]
            }
        })
        newFormData["persoanaContact"]["config"]["options"] = arrayOfContacts;
        setFormData(newFormData);
    }

    function getPropertyDataValuesToFormFields() {
        let newFormData = { ...formData };
        Object.keys(newFormData).forEach((prop, index) => {
            try {
                if (prop !== 'caracteristici' || prop !== 'specificatii' || prop !== "_id") {
                    newFormData[prop].value = props.propertyDetails[prop]
                }
            } catch (err) {
                console.log(err);
                console.log("FormData", newFormData);
                console.log("prop", prop);
            }

        });
        setFormData(newFormData);
    }

    // without caracteristici images and specificatii
    function getDataFromInput(e, inputName) {
        const updatedForm = { ...formData };
        const updatedElementForm = { ...updatedForm[inputName] };

        updatedElementForm.value = e.target.value;
        updatedForm[inputName] = updatedElementForm;

        setFormData(updatedForm);
    }

    /**
     * 
     * @param {Object} propFields 
     */
    function parsePropertyFields(propFields) {
        // daca este doar un tip de form
        // sa fac in asa fel incat sa selecteze tipul si dupa sa apara formul.
        let propertyForm;
        // <Form className={classes.Form} encType='multipart/form-data' action={`${url + requestUrls.ADD_PROPERTY}`} method="POST">
        // </Form>;
        /* sa creez doua input type hidden unde sa daug numarul de specificatii si caracteristici ca sa fie usor de parsat pe BE */
        propertyForm = (
            <Form className={classes.Form} onSubmit={(e) => { props.sendingData(e) }} >
                {
                    Object.keys(propFields).map((field, index) => {
                        return (<Input
                            key={index}
                            inputName={field}
                            inputValue={propFields[field]['value']}
                            inputConfig={propFields[field]['config']}
                            inputElement={propFields[field]['element']}
                            inputInfo={propFields[field]['info']}
                            change={getDataFromInput}
                        />)
                    })
                }
                {/* TODO lol got a quick fix here to send the property ID to BE but maybe I can send it in the request body without this  input */}
                {propertyEdit && <input type="hidden" id="property_id" name="property_id" value={propertyDetails._id} />}
                <Button type="submit" >{propertyEdit ? "Modifica proprietate" : "Adauga proprietate"}</Button>
            </Form>
        )
        return propertyForm;
    }

    return <div>
        {parsePropertyFields(formData)}
    </div>
}


export default PropertyForm;