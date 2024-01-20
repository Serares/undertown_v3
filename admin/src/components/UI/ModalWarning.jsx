import React, { useState, useEffect, useContext } from 'react';
import { Modal, Button } from 'react-bootstrap';

const ModalWarning = (props) => {

    let { modalInfo, abortAction, proceedAction } = props;
    let [show, setShow] = useState(false);

    const handleClose = () => setShow(false);

    useEffect(() => {
        setShow(modalInfo.showModal)
    })

    const modalFooter = () => {
        let contents = (
            <>
                <Button variant="secondary" onClick={() => { abortAction(); }}>
                    NU
          </Button>
                <Button variant="primary" onClick={proceedAction}>
                    Da
          </Button>
            </>
        )

        if (modalInfo.warning) {
            contents = (
                <>
                    <Button variant="secondary" onClick={() => { abortAction(); }}>
                        Ok
                    </Button>
                </>
            )
        }
        let modal_footer_elements = (
            <Modal.Footer>
                {contents}
            </Modal.Footer>
        );



        return modal_footer_elements;
    }

    return (
        <>
            <Modal show={show} onHide={handleClose} animation={false}>
                <Modal.Header closeButton>
                    <Modal.Title>{modalInfo.title}</Modal.Title>
                </Modal.Header>
                <Modal.Body>{modalInfo.message}</Modal.Body>
                {modalFooter()}
            </Modal>
        </>
    )
}


export default ModalWarning;