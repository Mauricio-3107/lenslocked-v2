package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Mauricio-3107/lenslocked-v2/models"
	"github.com/joho/godotenv"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "ffadde926c3b46"
	password = "8cbaa1078447b5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	es := models.NewEmailService(models.SMTPCongif{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err = es.ForgotPassword("mauro@m.com", "lenslocked.com/reset-ps?reseturl=abc123")
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent.")
}
