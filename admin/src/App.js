import React, { useState, useContext } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

import Layout from './hoc/Layout/Layout';
import AddProperty from './components/pages/AddProperty';
import Home from './components/pages/Home';
import Logout from './components/pages/Logout';
import Submissions from './components/pages/Submissions';
import ChangeProperty from './components/pages/ChangeProperty';
import ModalWarning from './components/UI/ModalWarning';

import { Route, Switch, withRouter } from 'react-router-dom';
import { GlobalStore } from './store/globalStore';

const App = props => {
  let modal_data = {
    message: "",
    title: "",
    showModal: false
  }

  const { setDeletedImages, domainUrl, endpoints, specificatiiData, caracteristiciData, deletedImages, setIsPropertyEdit } = useContext(GlobalStore);
  let [modalInfo, setModalInfo] = useState(modal_data);

  function sendDataToBackend(formData, isChangingProperty) {
    //TODO get rid of this postURL and add it to the global store
    let endpointUrl = domainUrl + "" + (isChangingProperty ? endpoints.CHANGE_PROPERTY : endpoints.ADD_PROPERTY);
    let METHOD = isChangingProperty ? "PUT" : "POST";

    console.log("Inputs", formData);
    let data = new FormData(formData);
    let valueCaracteristici = JSON.stringify(caracteristiciData);
    let valueSpecificatii = JSON.stringify(specificatiiData);
    data.append("caracteristici", valueCaracteristici);
    data.append("specificatii", valueSpecificatii);
    isChangingProperty && data.append('deletedImages', JSON.stringify(deletedImages));
    // for (let field of data) {
    //   console.log(field);
    // }
    fetch(endpointUrl, {
      method: METHOD,
      body: data
    })
    .then(res => {
        console.log(res);
        if (res.status != 200 && res.status != 201) {
          throw new Error('Failed request');
        }
        return res.json();
      })
      .then(responseData => {
        setDeletedImages([]);
        // in case property modified
        if (isChangingProperty) {
          setIsPropertyEdit(false);
        }
        setModalInfo({
          message: "Succes",
          title: "Proprietatea modificata/adaugata cu succes",
          showModal: true,
          warning: true
        })
        console.log(responseData);
      })
      .catch(err => {
        console.log(err);
        setModalInfo({
          message: "Eroare",
          title: "Eroare la adaugarea/modificarea proprietatii",
          showModal: true,
          warning: true
        })
      })
  }

  function closeModal() {
    setModalInfo(
      {
        message: "",
        title: "",
        showModal: false,
        warning: false
      });
  }

  let routes = (
    <Switch>
      <Route path="/admin/add" render={(p) => <AddProperty
        {...p}
        sendDataToBackend={sendDataToBackend}
      />} />
      <Route path="/admin/change/" render={(p) => <ChangeProperty
        {...p}
        sendDataToBackend={sendDataToBackend}
      />
      } />
      <Route path="/admin/submissions" component={Submissions} />
      <Route path="/admin/logout" component={Logout} />
      <Route exact path="/admin" render={(p) => <Home
        {...p}
      />} />
    </Switch>
  );

  return (
    <Layout>
      <ModalWarning
        abortAction={closeModal}
        proceedAction={() => { return false; }}
        modalInfo={modalInfo}
      />
      {routes}

    </Layout>)
};

export default withRouter(App);
