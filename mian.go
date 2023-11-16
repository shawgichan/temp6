package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "myprojectdb"
)

type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	PhoneNumber       string    `json:"phone_number"`
	OTP               string    `json:"otp"`
	OTPExpirationTime time.Time `json:"otp_expiration_time"`
}

func connectDB() (*pgx.Conn, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func createUser(conn *pgx.Conn, name string, phoneNumber string) error {
	query := `
    INSERT INTO users (name, phone_number)
    VALUES ($1, $2)
  `

	_, err := conn.Exec(context.Background(), query, name, phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func generateOTP(conn *pgx.Conn, phoneNumber string) (string, error) {
	otp := fmt.Sprintf("%04d", rand.Int31n(10000))

	expirationTime := time.Now().Add(1 * time.Minute)

	query := `
    UPDATE users 
    SET otp=$1, otp_expiration_time=$2
    WHERE phone_number=$3
  `

	_, err := conn.Exec(context.Background(), query, otp, expirationTime, phoneNumber)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func verifyOTP(conn *pgx.Conn, phoneNumber string, otp string) (bool, error) {
	query := `
    SELECT otp, otp_expiration_time 
    FROM users 
    WHERE phone_number = $1
  `

	var dbOTP string
	var expirationTime time.Time

	err := conn.QueryRow(context.Background(), query, phoneNumber).Scan(&dbOTP, &expirationTime)
	if err != nil {
		return false, err
	}

	if dbOTP != otp {
		return false, nil
	}

	if time.Now().After(expirationTime) {
		return false, nil
	}

	return true, nil
}

func main() {
	router := gin.Default()

	conn, err := connectDB()
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	router.POST("/api/users", func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := createUser(conn, user.Name, user.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully!"})
	})

	router.POST("/api/users/generateotp", func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		otp, err := generateOTP(conn, user.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"otp": otp})
	})

	router.POST("/api/users/verifyotp", func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		valid, err := verifyOTP(conn, user.PhoneNumber, user.OTP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully!"})
	})

	router.Run()
}
