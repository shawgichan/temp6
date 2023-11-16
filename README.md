Golang REST API Example



This is an example REST API built with Golang using Gin and PostgreSQL.

----------------------------------------------

Overview
The API implements basic user management functionality for:

----------------------------------------------

Creating a user
Generating OTP for a user
Verifying OTP for a user
It uses:

Gin for routing
pgx for PostgreSQL database access
PostgreSQL for data storage
Usage
Prerequisites
Golang
PostgreSQL

----------------------------------------------

Getting started
Clone the repo
Create the PostgreSQL database and table
Update the DB connection config in main.go
Run the server:
Copy code

$ go run main.go
The server will start on port 8080 by default.

----------------------------------------------

The API exposes the following endpoints:



POST /api/users
POST /api/users/generateotp
POST /api/users/verifyotp

See main.go for endpoint details.
