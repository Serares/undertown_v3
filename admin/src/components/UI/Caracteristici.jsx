import React, { useState, useEffect, useContext } from 'react';
import { Button, Form, Col } from 'react-bootstrap';

import { GlobalStore } from '../../store/globalStore';

const Caracteristici = (props) => {
    let { caracteristiciData, setCaracteristiciData } = useContext(GlobalStore);

    //caracteristici_data comes from Input as a value
    let { caracteristici_data } = props;

    useEffect(() => {
        if (caracteristici_data) {
            setCaracteristiciData(caracteristici_data);
        }
    }, [caracteristici_data]);

    function createCaracteristiciFields() {
        let newCaracteristiciData = [...caracteristiciData];
        let newCaracteristiciObject = { key: "", value: "" };
        newCaracteristiciData.push(newCaracteristiciObject);
        setCaracteristiciData(newCaracteristiciData);
    }

    function generateInputs() {

        let caracteristiciItems = caracteristiciData.map((caracteristica, index) => {
            return <li key={index}>
                <Form.Control as={"input"} onChange={(e) => { changeCaracteristicaItem(e, index, "key") }} type={"text"} placeholder={"Ex: Etaj"} id={"caracteristica-key-" + index} value={caracteristica.key} />
                <Form.Control as={"input"} onChange={(e) => { changeCaracteristicaItem(e, index, "value") }} type={"text"} placeholder={"Ex: 4/5"} id={"caracteristica-value-" + index} value={caracteristica.value} />
            </li>
        }
        );
        return caracteristiciItems;
    }

    function changeCaracteristicaItem(e, index, caracteristica_identificator) {
        console.log(e.target.value, index, caracteristica_identificator);
        let newCaracteristiciData = [...caracteristiciData];
        newCaracteristiciData[index][caracteristica_identificator] = e.target.value;
        setCaracteristiciData(newCaracteristiciData);
    }

    return <React.Fragment>
        <Button onClick={(e) => { e.preventDefault(); createCaracteristiciFields(); console.log("clicked add fields") }}>Adauga camp nou</Button>
        <Form.Row>
            <Form.Group as={Col}>
                <Form.Label>{"CARACTERISTICI"}</Form.Label>
                <ul className={"list-" + "CARACTERISTICI"} id="caracteristici-ul">
                    {generateInputs()}
                </ul>
            </Form.Group>
        </Form.Row>
    </React.Fragment>
}

export default Caracteristici;