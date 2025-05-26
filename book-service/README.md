---
title: Book Service
description: Implementing book service part of libray management backend app.
---

## Start

1. Build and start the containers:
    ```sh
    docker compose up --build
    ```

1. The application should now be running and accessible at `http://localhost:3000`.
   
## Endpoints

- `GET /api/v1/books`: Retrieves a list of all books.
  ```sh
  curl -X GET http://localhost:3000/api/v1/books
  ```

- `POST /api/v1/books`: Adds a new book to the collection.
  ```sh
  curl -X POST http://localhost:3000/api/v1/books \
       -H "Content-Type: application/json" \
       -d '{"title":"Title"}'
  ```
