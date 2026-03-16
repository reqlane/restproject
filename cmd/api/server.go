package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "restproject/internal/api/middlewares"
	"restproject/internal/api/router"
)

func main() {

	port := 3000
	cert := "cert.pem"
	key := "key.pem"
	router := router.Router()

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
	secureMux := mw.SecurityHeaders(router)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Handler: mux,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
