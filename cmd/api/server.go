package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restproject/internal/api/middlewares"
	"restproject/internal/api/router"
	"restproject/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error:", err)
		return
	}

	db, err := db.ConnectDb()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer db.Close()

	port := os.Getenv("SERVER_PORT")
	cert := "cert.pem"
	key := "key.pem"

	app := router.NewApp(db)
	mux := app.Router()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppConfig := mw.HPPConfig{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// secureMux := utils.ApplyMiddlewares(router,
	// 	mw.ResponseTime,
	// 	rl.RateLimit,
	// 	mw.Cors,
	// 	mw.SecurityHeaders,
	// 	mw.Hpp(hppConfig),
	// 	mw.Compression,
	// )
	secureMux := mw.SecurityHeaders(mux)

	server := http.Server{
		Addr: ":" + port,
		// Handler: mux,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
