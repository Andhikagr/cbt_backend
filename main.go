package main

import (
	"cbt-backend/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	
	fmt.Println("JWT Secret:", os.Getenv("JWT_SECRET"))



	
	r := routes.SetupRoutes()

	
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
