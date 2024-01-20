import React, { useState, useEffect } from 'react';
import { Map, TileLayer, Marker, Popup } from 'react-leaflet';
import classes from './Map.module.css';

const LeafletMap = (props) => {
    let { data } = props;

    let [lat, setLat] = useState(0);
    let [lng, setLng] = useState(0);
    let [zoom, setZoom] = useState(11);
    let [hasLocation, setHasLocation] = useState(false);

    useEffect(() => {
        try {

            setLat(data["lat"]);
            setLng(data["lng"]);
            setHasLocation(true);
        } catch (err) {
            console.log("Can't get data now");
        }
    }, [data])

    const handleClick = (e) => {

        console.log(e);
        setHasLocation(true);
        setLat(e["latlng"]["lat"]);
        setLng(e["latlng"]["lng"]);
        setZoom(17);
    }

    let position = [lat, lng];
    const createLocationValue = () => {
        let locObj = {
            lat: lat,
            lng: lng
        }

        return JSON.stringify(locObj);
    }
    const marker = (
        hasLocation ?
            <Marker position={position}>
                <Popup>Locația aleasă</Popup>
            </Marker> : null
    );

    return (<div className={classes.MapStyle}>
        <input type="hidden" name="location_coordonates" value={createLocationValue()} />
        <Map center={position} zoom={zoom} onclick={handleClick} >
            <TileLayer
                attribution='Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="https://www.mapbox.com/">Mapbox</a>'
                url="https://api.mapbox.com/styles/v1/empten/ck9qwupj76bgi1ipde38gjsmg/tiles/256/{z}/{x}/{y}@2x?access_token=pk.eyJ1IjoiZW1wdGVuIiwiYSI6ImNrOXF3eWh0azBvbXkzbHFjMWNoY2x1NzIifQ.xolL5C6kDqZvZzRFoCmdMg"
            />
            {marker}
        </Map>
    </div>)
}

export default LeafletMap;