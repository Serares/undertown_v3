import React from 'react';
import ReactDOM from 'react-dom';

import './index.css';
import App from './App';
import { BrowserRouter } from 'react-router-dom';
import { GlobalProvider } from './store/globalStore';

ReactDOM.render(
    <GlobalProvider>
        <BrowserRouter>
            <App />
        </BrowserRouter>
    </GlobalProvider>
    , document.getElementById('root'));
