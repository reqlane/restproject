package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restproject/internal/api/middlewares"
	"restproject/internal/api/routers"
	"restproject/internal/db"
	"time"

	"github.com/joho/godotenv"
)

// embedding .env into binary (Only for development)

//go:embed .env
var envFile embed.FS

func loadEnvFromEmbeddedFile() {
	content, err := envFile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	tempFile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating temp .env file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(content)
	if err != nil {
		log.Fatalf("Error writing to temp .env file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		log.Fatalf("Error closing temp .env file: %v", err)
	}

	err = godotenv.Load(tempFile.Name())
	if err != nil {
		log.Fatalf("Error loading temp .env file: %v", err)
	}
}

func main() {
	// Only in production, for running source code
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("Error:", err)
	// 	return
	// }

	loadEnvFromEmbeddedFile()

	db, err := db.ConnectDb()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer db.Close()

	port := os.Getenv("SERVER_PORT")
	// cert := "cert.pem"
	// key := "key.pem"
	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	app := routers.NewApp(db)
	mux := app.Router()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppConfig := mw.HPPConfig{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		WhiteList:                   []string{"sortBy", "sortOrder", "first_name", "last_name", "email", "class", "subject", "username"},
	}

	secureMux := mw.ApplyMiddlewares(mux,
		mw.ResponseTime,
		rl.RateLimit,
		mw.Cors,
		mw.SecurityHeaders,
		mw.WithPathsExcluded(mw.JWTMiddleware, "/execs/login", "/execs/forgotpassword", "/execs/resetpassword/reset"),
		mw.Hpp(hppConfig),
		mw.XSSMiddleware,
		mw.Compression,
	)

	server := http.Server{
		Addr:      ":" + port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
