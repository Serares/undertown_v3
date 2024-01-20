import React from 'react';
import { Button, Jumbotron } from 'react-bootstrap';

function Logout(props) {

    function postLogout() {
        fetch('/logout', {
            method: "POST",

        })
            .then(res => {
                console.log(res);
                window.location.replace(window.location.origin + "/login");
            })
            .catch(err => {
                console.log(err);
            })
    }

    function redirectToAdmin() {
        props.history.push('/admin');
    }

    return (
        <div className="logout">
            <Jumbotron>
                Sigur vrei sa te deloghezi?
            <Button onClick={postLogout} variant="primary">Da</Button>
                <Button onClick={redirectToAdmin} variant="warning">Nu, mai am treaba</Button>
            </Jumbotron>
        </div>
    )

}

export default Logout;