package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := 3000

	fmt.Println("Server is running on port:", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
