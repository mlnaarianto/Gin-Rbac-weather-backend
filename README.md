# Go Backend Application

A robust, convention-over-configuration backend application built with Go (Golang). This project structure is designed to mimic modern MVC/layered monolithic frameworks for high scalability and clean separation of concerns.

---

## 🚀 Getting Started

Follow these steps to get your development environment up and running.

### Prerequisites
Ensure you have the following software installed on your system:
*   **Go** (Version 1.20+) -> [Download Go](https://go.dev)
*   **Database Engine** (e.g., PostgreSQL, MySQL, or SQLite)

### 1. Installation
Clone the repository and download all the required project dependencies:
```bash
go mod tidy
```

### 2. Database Setup & Migrations
Configure your database connection strings inside your environment configuration, then populate your initial data using the built-in seeders:
```bash
# Run the database seeder to populate default roles and users
go run database/seeder.go
```

### 3. Running the Server
To start the local development server, run the following command in the root directory:
```bash
go run main.go
```
The server will boot up and start listening for HTTP requests (typically on `http://localhost:8080`).

---

## 📂 Project Architecture & Conventions

This application follows a strict directory layout to isolate responsibilities, heavily inspired by modern web framework design patterns:

### 🧩 Core Layers
*   **`main.go`**
    The main entry point of the application. It initializes configuration settings, connects to the database, plugs in middlewares, and boots the HTTP server.
*   **`routes/`**
    Acts as the router dispatching system (similar to `config/routes.rb` in Rails). Defines all application endpoints and maps HTTP verbs to their respective route handlers.
*   **`models/`**
    Contains the application data schemas and business logic representations. These structures map directly to your database tables.
*   **`middlewares/`**
    Global and route-specific interceptors. Handles cross-cutting concerns like Cross-Origin Resource Sharing (CORS) and authentication tokens.

### 🗄️ Database Management (`database/` & `db/`)
*   Configuration scripts for initializing connection pools.
*   Dedicated seeders (`seeder.go`, `user_seeder.go`, etc.) to programmatically inject initial or mock data into the system database for streamlined onboarding.

---

## 🛠️ Production Build

To compile the application into a single, production-ready executable binary:

```bash
# Build the binary executable
go build -o app main.go

# Execute the production server
./app
```