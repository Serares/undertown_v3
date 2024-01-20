import React, { useState, useContext, useEffect } from 'react';
import { Button, Form, Col } from 'react-bootstrap';
import styles from './ImagesInputs.module.css';
import { GlobalStore } from '../../store/globalStore';


const ImagesInputs = (props) => {
    let { isPropertyEdit, setDeletedImages, deletedImages } = useContext(GlobalStore);
    let { images_data, config } = props;
    let [images, setImages] = useState([]);

    //register deleted images only if isEditing is true
    //send to BE images that he deleted and all the images that are uploaded
    // replace the array in Mongo with the array that holds all the new Images
    // delete the images on S3 from the deleted images array

    useEffect(() => {
        setImages(images_data);
    }, [images_data])

    function addImage() {
        let newImages = [...images];
        newImages.push("");
        setImages(newImages);
    }

    function onImageChange(e, index) {
        let newImages = [...images];
        if (e.target.files && e.target.files[0]) {
            newImages[index] = URL.createObjectURL(e.target.files[0]);
        }
        setImages(newImages);
    }

    function renderUploadedImage(src, index) {
        let imageSrc = src;
        if (images_data) {
            imageSrc = (images_data.findIndex(img => img === src)) > -1 ? "/" + src : src;
        }

        return <img onClick={(e) => { deleteImage(src, index) }} className={styles.ImageElement} src={imageSrc} />
    }

    function deleteImage(src, index) {
        //add deleted in the deletedImages array
        let newDeletedImages = [...deletedImages];
        newDeletedImages.push(src);

        let newImages = [...images];
        console.log("Image index", newImages);
        newImages = newImages.filter((element, ind) => ind !== index);
        setDeletedImages(newDeletedImages);
        setImages(newImages);
    }

    function renderImagesInput() {
        if (!images) {
            return <Form.Label>---Nu sunt imagini adaugate</Form.Label>
        }

        let imagesInputs = images.map((image, index) => {
            return <Form.Group key={index} id={"image" + index}>
                <Form.Control {...config} onChange={(e) => { onImageChange(e, index) }} />
                {renderUploadedImage(image, index)}
            </Form.Group>
        })

        return imagesInputs;
    }


    return <>
        <Button onClick={addImage}>Adauga Imagine (max: 10)</Button>

        <div className={styles.ImagesInputs}>
            <Form.Label>Adauga Imaginile (Prima imagine adaugata o sa fie si la thumbnail)</Form.Label>
            <Form.Label>Poti sa stergi imaginea daca apesi pe ea</Form.Label>
            <div className={styles.ImagesInputsWrapper}>
                {renderImagesInput()}
            </div>
        </div>
    </>
}

export default ImagesInputs;
