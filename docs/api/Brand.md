# API Documentation: Brands

This document provides details about the API endpoints for managing brands.

## Brand Endpoints

| Method | Endpoint            | Description            | Request Body (JSON)                               |
| :--- | :------------------ | :--------------------- | :------------------------------------------------ |
| POST | `/api/v1/brands`    | Create a new brand     | `Brand` object (see structure below)              |
| GET  | `/api/v1/brands`    | Get all brands         | -                                                 |
| GET  | `/api/v1/brands/:id`| Get a brand by ID      | -                                                 |
| PUT  | `/api/v1/brands/:id`| Update a brand by ID   | `Brand` object (fields to be updated)             |
| DELETE | `/api/v1/brands/:id`| Delete a brand by ID   | -                                                 |

### `Brand` Object Structure

```json
{
  "name": "string",
  "imageUrl": "string"
}
```
