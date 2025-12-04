# Backend Application

A Go application built with Echo framework and PostgreSQL database.

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Configure environment variables:
   - Copy `.env.example` to `.env`
   - Update the `DATABASE_URL` with your database credentials

3. Run the application:
```bash
go run .
```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/example` - Get all examples
- `POST /api/v1/example` - Create a new example
- `GET /api/v1/example/:id` - Get example by ID
- `PUT /api/v1/example/:id` - Update example by ID
- `DELETE /api/v1/example/:id` - Delete example by ID

## Database Schema

Update the database schema in the `DATABASE_URL` connection string in your `.env` file.
