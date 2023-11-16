Here is one way to update the README to cover the sqlc usage:

Golang REST API Example
This is an example REST API built with Golang using Gin and PostgreSQL.

Overview
The API implements basic user management functionality for:

Creating a user
Generating OTP for a user
Verifying OTP for a user
It uses:

Gin for routing
pgx for PostgreSQL database access
PostgreSQL for data storage
I initially intended to also integrate sqlc for type-safe queries but faced issues setting it up on my Windows environment.

SQL Queries
The SQL queries for the main functions are defined in queries.sql:
-- Create user
INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING *;

-- Generate OTP 
UPDATE users SET otp = $1, otp_expiration = $2 WHERE id = $3;

-- Verify OTP
SELECT * FROM users WHERE phone_number = $1 AND otp = $2 AND otp_expiration > NOW();

These would be consumed by sqlc to generate type-safe Go code.

Usage
Prerequisites:

Golang
PostgreSQL
Getting started:

Clone the repo
Create the PostgreSQL database and table
Update the DB connection config in main.go
Run the server:

go run main.go

The server will start on port 8080 by default.

API Endpoints
The API exposes the following endpoints:

POST /api/users
POST /api/users/generateotp
POST /api/users/verifyotp
See main.go for endpoint details.
