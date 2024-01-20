import React from 'react';

import { Jumbotron, Navbar, Nav } from 'react-bootstrap';
import { NavLink } from 'react-router-dom';

import classes from './Layout.module.css';
// enclosing other components

const Layout = (props) => {

    const navLinksNames = {
        'Admin': '',
        'Adauga Proprietate': 'add',
        'Asteptare aprobare': 'submissions',
        'Logout': 'logout'
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
        <Jumbotron style={{ backgroundColor: '#fff' }} >
            <Navbar collapseOnSelect expand="lg" bg="dark" variant="dark">
                <Navbar.Brand>Admin</Navbar.Brand>
                <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                <Navbar.Collapse id="responsive-navbar-nav">
                    <Nav className="mr-auto">
                        {renderNavLinks()}
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
            {props.children}
        </Jumbotron>

    )

}

export default Layout;