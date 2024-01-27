import React, { useState, useEffect } from 'react';
import { Accordion } from 'react-bootstrap';
import PropertyCard from './PropertyCard';
import classes from './PropertiesTable.module.css';

/**
 * this will hold an array of PropertyCard
 * @param {*} props 
 */
const PropertiesTable = (props) => {

    let { properties, redirectToChangeProperty, deleteProperty } = props;

    useEffect(() => {
    }, [])

    function createPropertiesRows() {
        if (!properties) {
            return "Problem fetching properties"
        }

        if (properties.length < 1) {
            return "There are no properties added"
        }

        const elements = properties.map((elem, index) => {
            return <PropertyCard
                key={elem._id}
                property={elem}
                redirectToChangeProperty={redirectToChangeProperty}
                propertyIndex={index}
                deleteProperty={deleteProperty}
            />
        })
        return elements;
    }

    return (
        <div className={classes.PropertiesTable}>
            <Accordion>
                {createPropertiesRows()}
            </Accordion>
        </div>
    )

}

export default PropertiesTable;