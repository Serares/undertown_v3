import React, { useContext } from 'react';
import { Button } from 'react-bootstrap';
import { Carousel } from 'react-responsive-carousel';
import './carousel.css';

import classes from './PropertyCard.module.css';
import { GlobalStore } from '../../store/globalStore';

const PropertyCard = (props) => {
    let { property, redirectToChangeProperty, propertyIndex, deleteProperty } = props;
    let images = property.imagini;
    let { setPropertyDetails } = useContext(GlobalStore);
    let tipuriProprietate = [{ display: "Casa", value: "1" }, { display: "Apartament", value: "2" }, { display: "Garsoniera", value: "3" }, { display: "Vila", value: "4" }, { display: "Teren", value: "5" }, { display: "Birouri", value: "6" }, { display: "Proprietate Speciala", value: "7" }];
    let statusProprietate = [{ display: "Vanzare", value: "1" }, { display: "Inchiriere", value: "2" }];
    // let [loadedImages, setLoadedImages] = useState([]);
    // let [imagesLoaded, setIsImagesLoaded] = useState(false);

    //TODO make a image loading functionality
    // let loadImage = (imageSrc, index) => {

    //     return new Promise((resolve, reject) => {
    //         let imageElement = document.createElement('image');
    //         imageElement.src = imageSrc;
    //         imageElement.onload = () => {
    //             resolve(imageElement);
    //         }
    //         imageElement.onerror = () => {
    //             reject("Error loading image")
    //         }

    //     })

    // }

    // let imageElem = (src, index) => {
    //     let image;
    //     loadImage(src, index)
    //         .then(imageElement => {
    //             image = imageElement;
    //         })
    //         .catch(err => {
    //             console.log(err);
    //         })
    //     return image;
    // }

    // let loadAllImages = () => {
    //     let imagesNo = images.length;
    //     let imageElementsArray = [];
    //     for (let src of images) {

    //     }
    // }


    let carouselImages = (
        <div className={classes.CarouselContainer}>
            <Carousel>
                {images.map((elem, index) => {
                    return (
                        <div className={"21321"} key={index} >
                            <img src={"/" + elem} alt={"image_" + index} />
                        </div>
                    )
                })}
            </Carousel>
        </div>
    );

    function displayPropertyCaracteristics() {

        let chars = property.caracteristici.map((caracteristica, index) => {
            return (<ul key={index} className={'caracterisitca' + index}>
                <li>{caracteristica.key}</li>
                <li>{caracteristica.value}</li>
            </ul>)
        })
        return chars;
    }

    function displayPropertySpecifications() {
        try {
            return property.specificatii.map((spec, index) => {

                return (
                    <div key={"spec-title" + index}>
                        <p>{spec.name}</p>
                        <ul>
                            {spec.specs.map((spec_item, index) => {
                                return <li key={spec_item + " " + index}>
                                    {spec_item}
                                </li>
                            })}
                        </ul>
                    </div>
                )
            })
        } catch (err) {
            console.log("Eroare", err, property);
            return "";
        }

    }

    const changeHandler = () => {
        setPropertyDetails(property);
        redirectToChangeProperty(property._id, propertyIndex);
    }

    const deleteHandler = () => {
        deleteProperty(property._id, propertyIndex);
    }

    const displayPropertyType = function (propertyType) {
        let displayValue = "";
        tipuriProprietate.forEach(type => {
            if (+type.value === propertyType) {
                displayValue = type.display;
            }
        });

        return displayValue;
    }

    const statusProperty = (propertyStatus) => {
        let displayValue = "";
        statusProprietate.forEach(status => {
            if (+status.value === propertyStatus) {
                displayValue = status.display;
            }
        });

        return displayValue;
    }
    return (
        <div className={classes.PropertyCard}>
            {carouselImages}
            <div className={classes.PropertyDetails}>
                <div className={classes.Titlu}>Titlu proprietate: {property.titlu}</div>
                <div className={classes.Adresa}>Adresa proprietate: {property.adresa}</div>
                <div className={classes.Descriere}>Descriere proprietate: {property.detalii}</div>
                <div className={classes.TipProprietate}>Tip proprietate: {displayPropertyType(property.tipProprietate)}</div>
                <div className={classes.Status}>Status prop: {statusProperty(property.status)}</div>
                <div className={classes.Pret}>Pret: {property.pret}€</div>
                <div className={classes.Suprafata}>Suprafata: {property.suprafata} mp</div>
                <div className={classes.Caracteristici}>Caracteristici: {displayPropertyCaracteristics()}</div>
                <div className={classes.Specificatii}>Specificatii: {displayPropertySpecifications()}</div>
                <div className={classes.Featured}>Este in featured: {property.featured === 0 ? "Nu" : "Da"}</div>
            </div>
            <div className={classes.PropertyEvents}>
                <Button variant="warning" onClick={changeHandler} >Modifică</Button>
                <Button variant="danger" onClick={deleteHandler} >Șterge</Button>
            </div>
        </div>
    )

}

export default PropertyCard;