import React, { useState, useEffect, useContext } from 'react';
import { Button, Form, Col } from 'react-bootstrap';

import { GlobalStore } from '../../store/globalStore';

const Specificatii = (props) => {
    let { specificatiiData, setSpecificatiiData } = useContext(GlobalStore);

    //specificatii_data comes from Input as a value
    let { specificatii_data, onchange } = props;

    useEffect(() => {
        if (specificatii_data) {
            setSpecificatiiData(specificatii_data);
        } else {
            setSpecificatiiData([{ name: "", specs: [] }]);
        }
    }, [specificatii_data])

    function parseSpecsDataToArray(data) {
        let newData = [];
        for (let spec of data) {
            let newSpec = {}
            newSpec[spec] = [...data[spec]];
            newData.push(newSpec);
        }
        return newData;
    }

    /**
     * 
     * @param {Event} e 
     * @param {Number} liIndex 
     * @param {String} spec_name 
     */
    function changeSpecsLi(e, liIndex, spec_name) {
        let newSpecsData = [...specificatiiData];
        let newSpec = newSpecsData.find(spec => spec.name === spec_name);
        let newSpecSpecs = newSpec.specs;
        newSpecSpecs[liIndex] = e.target.value;
        setSpecificatiiData(newSpecsData);
    }

    function incrementListItems(index) {
        let newSpecsData = [...specificatiiData];
        newSpecsData[index].specs.push("");
        setSpecificatiiData(newSpecsData);
    }

    /**
     * 
     * @param {String[]} specs_list 
     * @param {String} spec_name 
     */
    function generateListItems(specs_list, spec_name) {
        try {
            let listItems = specs_list.slice();
            let listItemsElements = listItems.map((item, index) => {
                return <Form.Control value={item} key={index} as={"input"} onChange={(e) => { changeSpecsLi(e, index, spec_name) }} type={"text"} placeholder={"Ex:Internet"}
                    id={"element-specificatie-" + index} />
            })
            return listItemsElements;
        } catch (err) {
            console.log("Specificatii error ", err);
            return "";
        }

    }

    function addList() {
        let newSpecsData = [...specificatiiData];
        let newSpecObject = { name: "", specs: [] };
        newSpecsData.push(newSpecObject);

        setSpecificatiiData(newSpecsData);
    }

    /**
     * 
     * @param {Event} e 
     * @param {String} listTitle 
     * @param {Number} index 
     */
    function changeListTitle(e, listTitle, index) {
        let newSpecsData = [...specificatiiData];
        if (listTitle !== e.target.value) {
            newSpecsData[index].name = e.target.value;
        }
        setSpecificatiiData(newSpecsData);
    }

    function generateLists() {
        if (!specificatiiData) {
            return [];
        }

        let elements = specificatiiData.map((specificatie, index) => {
            return <Form.Row key={index}>
                <Form.Group as={Col}>
                    <Form.Control as="input" onChange={(e) => { changeListTitle(e, specificatie.name, index) }} value={specificatie.name} type="text" placeholder="Ex:Utilitati" id={"titlu-specificatii-" + index} />
                    <Button onClick={(e) => { incrementListItems(index) }}>Adauga specificatie pentru {specificatie.name}</Button>
                    <ul className={"list-" + "inputName"} id={"lista-specificatii-" + index}>
                        {generateListItems(specificatie.specs, specificatie.name)}
                    </ul>
                </Form.Group>
            </Form.Row>
        })
        return elements;
    }

    return (
        <div>
            <Button onClick={(e) => { addList() }}>Adauga lista specificatii</Button>
            {generateLists()}
        </div>
    );
}

export default Specificatii;