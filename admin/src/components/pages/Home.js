import React, { useState, useEffect, useContext } from 'react';
import fetchData from '../../service/fetchService';

import PropertiesTable from '../UI/PropertiesTable';
import LoadingIndicator from '../UI/LoadingIndicator';

import { GlobalStore } from '../../store/globalStore';
import ModalWarning from '../UI/ModalWarning';

const Home = (props) => {
    let modal_data = {
        message: "",
        title: "",
        showModal: false
    }

    let property_delete_data = {
        _id: "",
        index_number: ""
    }

    let { setProperties, domainUrl, endpoints } = useContext(GlobalStore);
    let [dbProperties, setDbProperties] = useState([]);

    let [isLoading, setIsLoading] = useState(true);
    let [modalInfo, setModalInfo] = useState(modal_data);
    let [propertyForDeletion, setPropertyForDeletion] = useState(property_delete_data);
    let [isComponentRendered, setIsComponentRendered] = useState(false);

    useEffect(() => {
        setIsComponentRendered(true);
    }, [])

    useEffect(() => {
        const dataUrl = domainUrl + '' + endpoints.GET_PROPERTIES;
        fetchData(dataUrl)
            .then(data => {
                setDbProperties(data.properties);
                setIsLoading(false);
            })
            .catch(err => {
                console.log(err);
                setModalInfo(
                    {
                        message: "Eroare de retea, nu se poate conecta la baza de date",
                        title: "Eroare",
                        showModal: true,
                        warning: true
                    })
            })
    }, [isComponentRendered]);

    function deleteProperty(propertyId, propertyIndex) {
        setPropertyForDeletion({
            _id: propertyId,
            index_number: propertyIndex
        });

        setModalInfo(
            {
                message: "Esti sigur ca vrei sa stergi proprietatea " + dbProperties[propertyIndex].titlu,
                title: "Stergi proprietatea?",
                showModal: true,
                warning: false
            })
    }

    function abortDeleteProperty() {
        setModalInfo(
            {
                message: "",
                title: "",
                showModal: false,
                warning: false
            });
        setPropertyForDeletion({});
    }

    function updatePropertiesAfterDeletion(propertyDetails) {
        let newProperties = [...dbProperties];

        newProperties = newProperties.filter(prop => prop._id !== propertyDetails._id);
        setProperties(newProperties);
        setDbProperties(newProperties);
    }

    function sendDeleteData() {

        fetch(domainUrl + '' + endpoints.DELETE_PROPERTY + '/' + propertyForDeletion._id, {
            method: "DELETE"
        })
            .then(res => {
                console.log(res.status);
                if (res.status !== 200 && res.status !== 201) {
                    throw new Error('Failed to delete property');
                }

                return res.json()
            })
            .then(data => {
                setModalInfo({
                    message: "Properietatea a fost stearsa din baza de date",
                    title: "Success",
                    showModal: true,
                    warning: true
                })
                updatePropertiesAfterDeletion(propertyForDeletion);
                setPropertyForDeletion({});
            })
            .catch(err => {
                console.log(err);
                setModalInfo({
                    message: "Erroare la stergerea proprietatii incearca din nou",
                    title: "Eroare",
                    showModal: true,
                    warning: true
                })
            })
    }

    const redirectToChangeProperty = (propertyID, index) => {
        props.history.push(props.match.path + '/change/?objectId=' + propertyID + '&propertyIndex=' + index);
    }

    return (
        <div>
            <ModalWarning
                modalInfo={modalInfo}
                proceedAction={sendDeleteData}
                abortAction={abortDeleteProperty}
            />
            Search property <input type="text" />
            {isLoading ?
                <LoadingIndicator /> :
                <PropertiesTable
                    properties={dbProperties}
                    redirectToChangeProperty={redirectToChangeProperty}
                    deleteProperty={deleteProperty}
                />}
        </div>
    )
}


export default Home;