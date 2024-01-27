import React from 'react';

import { Accordion } from 'react-bootstrap';
import { NavLink } from 'react-router-dom';
import {Navbar} from '../../components/UI/Navbar'
import classes from './Layout.module.css';
// enclosing other components

const Layout = (props) => {

    const navLinksNames = {
        'Admin': '',
        'Adauga Proprietate': 'add',
        'Asteptare aprobare': 'submissions',
        'Logout': 'logout',
        'Login' : 'login'
    }

    function renderNavLinks() {
        let navLinks = [];
        navLinks = Object.keys(navLinksNames).map((item, index) => {
            
            return (
                <li className={classes.NavigationItem} key={index}>
                    <NavLink exact activeClassName={classes.active} to={`/admin/${navLinksNames[item]}`}>
                        {item}
                    </NavLink>
                </li>)
        });

        return navLinks;
    }

    return (
        <Accordion style={{ backgroundColor: '#fff' }} >
            {props.children}
        </Accordion>

    )

}

export default Layout;