# Go Backend Service Template

A clean architecture Go backend service template using net/http, sqlx, and PostgreSQL.

## Project Structure

```
.
├── config
│   └──  database.ini     # Database configuration file
├── db
│   └── migrations        # Database migration files
├── internal
│   ├── config           # Configuration management
│   ├── entity           # Database entities
│   ├── handler          # HTTP handlers (controllers)
│   ├── model            # Data models
        ├── converter     # Converter functions
│   ├── repository       # Database operations
│   ├── router          # HTTP routing
│   └── service         # Business logic
│   └── utils           # Utility functions
├── go.mod              # Go module file
├── cmd                  # Command-line applications
│   └── main.go          # Main entry point
└── README.md           # Project documentation
```

## Prerequisites

- Go 1.21.5 or later
- PostgreSQL

## Setup

1. Create a PostgreSQL database:

```sql
CREATE DATABASE go_template;
```

2. Run the database migrations:

```sql
-- Execute the SQL files in the db/migrations directory
```

3. Start the server:

```bash
go run cmd/main.go
```

TODO:

- Add a Makefile for common tasks
- Add a Dockerfile for containerization
- Add a README.md for project documentation
- Unit tests
- CI/CD pipeline

The server will start on `localhost:8080`
