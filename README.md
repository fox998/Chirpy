## 🐦 Chirpy - A Structured Go HTTP Server

**Chirpy** is an in-depth **educational project** implementing a full-featured, foundational **HTTP server** in **Go (Golang)**. It demonstrates best practices for structuring a Go application and integrating professional data management tools.

This repository serves as a practical exercise and final project for the **Boot.dev** course: [**Learn HTTP Servers with Go**](https://www.boot.dev/courses/learn-http-servers-golang).

-----

### 🚀 Getting Started

To run the Chirpy server locally, you will need **Go** installed and a running **PostgreSQL** instance for database persistence.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/fox998/Chirpy.git
    cd Chirpy
    ```

2.  **Setup the Database (PostgreSQL):**

      * Set your database connection string as an environment variable (e.g., `DB_URL=...`).
      * Run migrations using the included `goose.sh` script to set up the schema.

3.  **Run the server:**

    ```bash
    go run .
    ```

The server will start and be ready to handle requests, serving the basic web interface from `index.html` and the `assets/` directory.

-----

### 🏗️ Architecture & Key Technologies

The project is structured to reflect professional Go development standards:

#### Application Structure

  * **`main.go`**: The application entry point, responsible for configuration, connection pooling, and setting up the router.
  * **`internal/`**: Contains all private application logic, including the HTTP handlers, business logic, and database access layer, promoting **separation of concerns**.
  * **`assets/`**: Hosts static files and the primary **`index.html`** frontend.

#### Data Persistence Tools

This server goes beyond basic examples by implementing a robust database layer using:

  * **PostgreSQL**: The chosen relational database (setup implied by `run_psql.sh`).
  * **Goose**: Used for **database migrations**, ensuring the schema is version-controlled and easily managed.
  * **SQLC**: Used to **generate type-safe Go code** directly from raw SQL queries, significantly reducing boilerplate and increasing data access reliability.

-----

### 📚 Core Functionality

The code demonstrates core concepts of building a persistent, authenticated server:

  * **RESTful API**: Provides endpoints for managing **"Chirps"** (posts) and **"Users"**.
  * **Authentication**: Implements secure **user registration**, **login**, and **token-based authorization** to protect sensitive API routes.
  * **Middleware**: Uses handler wrapping for common tasks like **logging** and **authentication checks** before requests reach the core handler logic.
  * **Data Handling**: Manages requests and responses using structured JSON data.