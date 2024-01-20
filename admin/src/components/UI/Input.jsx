import React, { useState, useEffect } from 'react';
import { Form, Col } from 'react-bootstrap';
import Specificatii from './Specificatii';
import Caracteristici from './Caracteristici';
import ImagesInputs from './ImagesInputs';
import LeafletMap from './Map';

const Input = (props) => {

    let value = props.inputValue;
    const inputName = props.inputName;
    const config = props.inputConfig;
    const element = props.inputElement;
    let [inputValue, setInputValue] = useState();
    let inputInfo = props.inputInfo;
    useEffect(() => {
        setInputValue(value)
    }, [value]);

    let input;
    switch (element) {
        case ('input'):
            input = (
                <Form.Row>
                    <Form.Group as={Col}>
                        <Form.Label>{inputName}</Form.Label>
                        <Form.Text className="text-muted">
                            {inputInfo}
                        </Form.Text>
                        <Form.Control as={element} onChange={(e) => { props.change(e, inputName) }} {...config} value={inputValue} />
                    </Form.Group>
                </Form.Row>
            );
            break;
        case ('select'):
            input = (
                <Form.Row>
                    <Form.Group as={Col}>
                        <Form.Label>{inputName}</Form.Label>
                        <Form.Text className="text-muted">
                            {inputInfo}
                        </Form.Text>
                        <Form.Control value={inputValue} as={element} onChange={(e) => { props.change(e, inputName) }} type={config.type} placeholder={config.placeholder} id={inputName} name={inputName}>
                            {config.options.map((option, index) => {
                                return <option key={option["value"]} value={option["value"]}>{option["display"]}</option>
                            })}
                        </Form.Control>
                    </Form.Group>
                </Form.Row>);
            break;
        case ('textarea'):
            input = (
                <Form.Row>
                    <Form.Group as={Col}>
                        <Form.Label>{inputName}</Form.Label>
                        <Form.Control as={element} onChange={(e) => { props.change(e, inputName) }} rows="3" {...config} value={inputValue} />
                    </Form.Group>
                </Form.Row>);
            break;
        case ('caracteristici'):
            input = (
                <Caracteristici caracteristici_data={inputValue} />
            );
            break;
        case ('specificatii'):
            input = <Specificatii specificatii_data={inputValue} onchange={props.change} />
            break;
        case ('inputFiles'):
            input = (
                <ImagesInputs images_data={inputValue} config={config} />
            );
            break;
        case ('location'):
            input = (
                <div>
                    <h2>Adaugă locația proprietății pe hartă</h2>
                    <LeafletMap data={props.inputValue} />
                </div>
            );
            break;
        default:
            input = "Default input";
            break;
    }

    return input;
}

export default Input;
