# API

## Overview

This API allows for managing property listings in a real estate application. It supports operations to add, update, retrieve, and delete property listings. Additionally, there are server-side rendered webpages for admin users.

## Authorization

All API requests require the use of a JWT token. This token must be included in the `Authorization` header of each request. The token is obtained through the API Gateway authorizer.

## API Endpoints

### Add Property

- **URL**: `/addProperty`
- **Method**: `POST`
- **Description**: Adds a new property listing.
- **Headers**:
  - `Authorization`: `Bearer YOUR_JWT_TOKEN`
- **Body**:
  ```json
  {
    "title": "string",
    "description": "string",
    "price": "number"
    // Other fields...
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**: `{ "message": "Property added successfully", "id": "property_id" }`

### Update Property

- **URL**: `/addProperty`
- **Method**: `PUT`
- **Description**: Updates an existing property listing.
- **Headers**:
  - `Authorization`: `Bearer YOUR_JWT_TOKEN`
- **Body**:
  ```json
  {
    "id": "property_id",
    "title": "string",
    "description": "string",
    "price": "number"
    // Other fields...
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**: `{ "message": "Property updated successfully" }`

### Delete Property

- **URL**: `/deleteProperty`
- **Method**: `DELETE`
- **Description**: Deletes a property listing.
- **Headers**:
  - `Authorization`: `Bearer YOUR_JWT_TOKEN`
- **Query Parameters**:
  - `id`: Property ID to delete
- **Success Response**:
  - **Code**: 200
  - **Content**: `{ "message": "Property deleted successfully" }`

### Get Property

- **URL**: `/getProperty`
- **Method**: `GET`
- **Description**: Retrieves details of a specific property.
- **Headers**:
  - `Authorization`: `Bearer YOUR_JWT_TOKEN`
- **Query Parameters**:
  - `id`: Property ID to retrieve
- **Success Response**:
  - **Code**: 200
  - **Content**: `{ "property": { /* Property Details */ } }`

### Get Properties

- **URL**: `/getProperties`
- **Method**: `GET`
- **Description**: Retrieves a list of all properties.
- **Headers**:
  - `Authorization`: `Bearer YOUR_JWT_TOKEN`
- **Success Response**:
  - **Code**: 200
  - **Content**: `{ "properties": [ /* Array of Properties */ ] }`

## Admin Webpages

### List Properties

- **URL**: `/list`
- **Method**: `GET`
- **Description**: Displays a list of all properties for admin users.

### Submit Property

- **URL**: `/submit`
- **Methods**: `GET`, `POST`
- **Description**:
  - `GET`: Displays a form to submit a new property.
  - `POST`: Submits a new property.

### Edit Property

- **URL**: `/edit/<property_title>?propertyId=<id>`
- **Methods**: `PUT`, `GET`
- **Description**:
  - `PUT`: Submits the updates for an existing property.
  - `GET`: Get's the edit property webpage

### Delete Property

- **URL**: `/delete`
- **Method**: `DELETE`
- **Description**: Deletes a property from the admin interface.

## Homepage Webpages

### List Properties

- **URL**: `/chirii`
- **Method**: `GET`
- **Description**: Displays the `RENT` `property_type` of properties.

### List Properties

- **URL**: `/vanzari`
- **Method**: `GET`
- **Description**: Displays a list of all properties for admin users.

### Homepage

- **URL**: `/`
- **Methods**: `GET`
- **Description**:
  - `GET`: Displays a list of []lite.ListFeaturedPropertiesRow list.
