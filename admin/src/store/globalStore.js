import React, { useState, useEffect } from 'react';
import fetchData from '../service/fetchService';

const GlobalStore = React.createContext();
const { Provider } = GlobalStore;

const GlobalProvider = ({ children }) => {
    const routes = {
        ADD_PROPERTY: "/admin/postAddProperty",
        DELETE_PROPERTY: "/admin/deleteProperty",
        MODIFY_PROPERTY: "/admin/modifyProperty",
        GET_PROPERTIES: "/admin/getAllProperties",
        CHANGE_PROPERTY: "/admin/changeProperty",
        GET_ADMIN_ACCOUNTS: "/admin/get_admin_accounts",
        GET_PROPERTY_FORM_FIELDS: "/admin/get_property_form_fields"
    };

    const hostUrl = window.location.protocol + "//" + window.location.host;
    const [domainUrl, setDomainUrl] = useState(hostUrl);
    const [endpoints, setEndpoints] = useState(routes);
    const [propertyDetails, setPropertyDetails] = useState({});
    let [properties, setProperties] = useState([]);
    let [specificatiiData, setSpecificatiiData] = useState([]);
    let [caracteristiciData, setCaracteristiciData] = useState([]);
    let [isPropertyEdit, setIsPropertyEdit] = useState(false);
    //getting administrator accounts to be used in submitProperty
    const [adminAccounts, setAdminAccounts] = useState([]);
    const [propertyFields, setPropertyFields] = useState({});
    // only used for propertyChange
    let [deletedImages, setDeletedImages] = useState([]);


    function getAdminAccounts() {
        fetchData(domainUrl + endpoints.GET_ADMIN_ACCOUNTS, {
            method: "GET"
        })
            .then(response => {
                console.log("accounts", response["accounts"]);
                setAdminAccounts(response["accounts"]);
            })
            .catch(error => {
                console.log("Error")
            })
    }

    function getPropertyFields() {
        let propertyFields = {};

        fetchData(domainUrl + endpoints.GET_PROPERTY_FORM_FIELDS,
            { method: "GET" })
            .then(response => {
                console.log(JSON.parse(response[0]["propertyFields"]));
                propertyFields = JSON.parse(response[0]["propertyFields"]);
                propertyFields["tipProprietate"]["config"]["options"] = JSON.parse(response[0]["tipuriProprietate"]);
                propertyFields["status"]["config"]["options"] = JSON.parse(response[0]["statusProprietate"]);
                setPropertyFields(propertyFields);
            })
            .catch(err => {
                console.log(err)
            })
    }

    useEffect(() => {
        console.log("Store mounted");
        getAdminAccounts();
        getPropertyFields();
    }, [])

    const state = {
        domainUrl,
        endpoints,
        propertyDetails,
        properties,
        specificatiiData,
        caracteristiciData,
        isPropertyEdit,
        deletedImages,
        adminAccounts,
        propertyFields
    };

    const actions = {
        setDomainUrl,
        setEndpoints,
        setPropertyDetails,
        setSpecificatiiData,
        setCaracteristiciData,
        setIsPropertyEdit,
        setDeletedImages,
        setProperties,
        setAdminAccounts
    }

    return (<Provider value={{ ...state, ...actions }}>{children}</Provider>)
}

export { GlobalProvider, GlobalStore }