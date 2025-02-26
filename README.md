# Simpler Products API

This is a Go-based RESTful API for managing products in a store. It provides endpoints for creating, reading, updating, and deleting products, along with basic authentication and pagination support. The API is built using the Gin web framework and interacts with a MySQL database for data storage.

## Features

* **Product CRUD operations:**
  * `GET /api/v1/products`: Retrieve a list of products with pagination support.
  * `GET /api/v1/products/:id`: Retrieve a specific product by its ID.
  * `POST /api/v1/products`: Create a new product.
  * `PUT /api/v1/products/:id`: Update an existing product.
  * `DELETE /api/v1/products/:id`: Delete a product.
* **Authentication:**
  * `JWT_SECRET_KEY` and `AUTH_ENABLED` environment variables control JWT authentication.
  * Product endpoints require a valid JWT token in the `Authorization` header when authentication is enabled.
* **Pagination:**
  * The `GET /api/v1/products` endpoint supports pagination using `limit` and `offset` query parameters.
* **Response handling:**
  * Centralized response and error handling middleware provides consistent and structured responses.
* **Logging:**
  * Uses `logrus` for structured logging.
* **Validation:**
  * Basic input validation is implemented using Gin's binding and validation features.
* **Testing:**
  * Unit tests are included for services, handlers, validators, and middleware.
* **Docker support:**
  * A `Dockerfile` and `docker-compose.yml` file are provided for easy containerization and deployment.

## Prerequisites

* **Go:** Make sure you have Go installed on your system. You can download it from the official website: [https://golang.org/dl/](https://golang.org/dl/)
* **MySQL:** You'll need a running MySQL server. You can install it locally or use a Docker container (see the "Running with Docker Compose" section below).
* **Docker (optional):** If you want to run the application using Docker, make sure you have Docker installed. You can download it from the official website: [https://www.docker.com/get-started](https://www.docker.com/get-started)
* **Docker Compose (optional):** If you want to use Docker Compose for easier setup, make sure you have it installed. You can find instructions on the official website: [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

## Getting Started

1. **Clone the repository:**

    ```bash
    git clone https://github.com/dimitrisniras/simpler-products.git
    cd simpler-products
    ```

2. **Set up environment variables:**

    * Create a `.env` file in the project root directory.

    * Add the following environment variables, replacing the placeholders with your actual values:

        ```yml
        PORT=8080
        LOG_LEVEL=debug # or 'trace', 'info', 'warn', 'error', 'release'
        DB_USER=your_db_user
        DB_PASSWORD=your_db_password
        DB_HOST=localhost # or the hostname/IP of your MySQL server
        DB_PORT=3306
        DB_NAME=your_db_name
        JWT_SECRET_KEY=your_strong_secret_key
        AUTH_ENABLED=true # or 'false' to disable authentication
        ```

3. **Create the database and table:**

    * Connect to your MySQL server using a tool like the MySQL command-line client or MySQL Workbench.

    * Create the database specified in your `.env` file (e.g., `your_db_name`).

    * Execute the following SQL query to create the `products` table:

        ```sql
        CREATE TABLE Products (
            id VARCHAR(255) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description VARCHAR(255) NOT NULL,
            price DECIMAL(10, 2) NOT NULL
        );
        ```

4. **Install dependencies:**

    ```bash
    make install-deps
    ```

## Running the Application

### Locally

1. **Build the executable:**

    ```bash
    make build
    ```

2. **Run the application:**

    ```bash
    make run
    ```

    The API will be accessible at `http://localhost:8080`.

### With Docker

1. **Build and run the containers:**

    ```bash
    docker-compose up --build
    ```

    This will build the Go application image, pull the MySQL image, and start both containers. The API will be accessible at `http://localhost:8080`.
    Remember if you're a Mac OS user you'll need to specify the DB_HOST as `docker.for.mac.localhost`

## Testing

### Unit Tests

* **Run unit tests:**

    ```bash
    make test
    ```

    This will execute all the unit tests in the `tests` directory.

## Authentication

This API utilizes JWT (JSON Web Token) for authentication. You'll need to generate an RSA key pair (private and public keys) and set the public key as an environment variable.

### Generating Keys

You can generate an RSA key pair using OpenSSL:

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem
```

* The `private.pem` file contains your private key, which should be kept secure and not shared.
* The `public.pem` file contains your public key, which will be used by the API to verify JWTs.

### Setting up the Public Key

1. **Encode the public key to Base64:**

    ```bash
    export JWT_PUBLIC_KEY=$(cat public.pem | base64 -w 0)
    ```

2. **Set the environment variable:**

* **Locally:** Include the following line in your `.env` file, replacing `your_base64_encoded_public_key` with the actual Base64-encoded value of your public key:

    ```bash
    JWT_PUBLIC_KEY=your_base64_encoded_public_key
    ```

* **Docker Compose:** Add the following environment variable to the app service in your `docker-compose.yml` file:

    ```yaml
    environment:
    - JWT_PUBLIC_KEY=your_base64_encoded_public_key
    # ... other environment variables
    ```

## API Endpoints

* **`GET /api/ping`**

  * A simple health check endpoint that returns 200 OK.
  * Does not require authentication.

  * **Success Response:**

    ```json
    {
        "status": 200,
    }
    ```

* **`GET /api/v1/products`**

  * Retrieves a list of products.
  * Supports pagination using limit and offset query parameters.
  * Requires authentication (when AUTH_ENABLED is true).

  * **Success Response (with pagination):**

    ```json
    {
        "status": 200,
        "data": [
            {
                "id": "uuid1",
                "name": "Product A",
                "description": "Description of Product A",
                "price": 10.99
            },
            // ... other products
        ],
        "pagination": {
            "limit": 10,
            "offset": 0,
            "total": 25 
        }
    }
    ```

    * **Error Response (e.g., Invalid parameters):**

    ```json
    {
        "status": 400,
        "errors": [
            {
                "message": "Invalid limit parameter"
            }
        ]
    }
    ```

* **`GET /api/v1/products/:id`**

  * Retrieves a specific product by its ID.
  * Requires authentication.

  * **Success Response:**

    ```json
    {
        "status": 200,
        "data": [{
            "id": "uuid1",
            "name": "Product A",
            "description": "Description of Product A",
            "price": 10.99
        }]
    }
    ```

    * **Error Response (e.g., Product not found):**

    ```json
    {
        "status": 404,
        "errors": [
            {
                "message": "product not found"
            }
        ]
    }
    ```

* **`POST /api/v1/products`**

  * Creates a new product.
  * Requires authentication.
  
  * **Success Response:**

    ```json
    {
        "status": 201,
        "data": [{
            "id": "uuid1",
            "name": "New Product",
            "description": "This is a new product",
            "price": 19.99
        }]
    }
    ```

  * **Error Response (e.g., Validation errors):**

    ```json
    {
        "status": 400,
        "errors": [
            {
                "message": "Name is required"
            },
            {
                "message": "Price must be greater than 0"
            }
        ]
    }
    ```

* **`PUT /api/v1/products/:id`**

  * Updates an existing product.
  * Requires authentication.

  * **Success Response:**

    ```json
    {
        "status": 200,
        "data": [{
            "id": "uuid1",
            "name": "Updated Product",
            "description": "Updated description",
            "price": 24.95
        }]
    }
    ```

    * **Error Response (e.g., Product not found):**

    ```json
    {
        "status": 404,
        "errors": [
            {
                "message": "product not found"
            }
        ]
    }
    ```

* **`DELETE /api/v1/products/:id`**

  * Deletes a product.
  * Requires authentication.

  * **Success Response:**

    ```json
    {}
    ```

  * **Error Response (e.g., Product not found):**

    ```json
    {
        "status": 404,
        "errors": [
            {
                "message": "product not found"
            }
        ]
    }
    ```

## Examples

### Creating a Product

```bash
curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Bearer your_jwt_token" \
-d '{"name": "New Product", "description": "This is a new product", "price": 19.99}' \
http://localhost:8080/api/v1/products
```

### Retrieving Products

```bash
curl -H "Authorization: Bearer your_jwt_token" \
http://localhost:8080/api/v1/products?limit=5&offset=0
```

Remember to replace `your_jwt_token` with an actual valid JWT token if authentication is enabled.
