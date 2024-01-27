import React, { useState, useContext } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import {Navbar} from './components/UI/Navbar'
import Layout from './hoc/Layout/Layout';
import { Route, Routes } from 'react-router-dom';
import Login from './components/pages/Login';
import SubmitProperty from './components/pages/SubmitProperty';

const App = props => {
  
  let routes = (
    <Routes>
      <Route path="/login" element={<Login />}/>
      <Route path="/add" element={<SubmitProperty />}/>
    </Routes>
  );

  return (
    <div>
      <Navbar />
      {routes}
    </div>)
};

export default App;
