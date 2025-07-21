package main

import (
	"cbt-backend/routes"
	"fmt"      // Untuk mencetak output ke terminal
	"log"      // // Untuk mencetak log dan menangani error log
	"net/http" // Untuk menjalankan server HTTP
	"os"       // Untuk membaca environment variable dari OS

	"github.com/joho/godotenv" // Library untuk membaca file .env
)

func main() {
	// Load .env dulu
	// Memuat file .env agar isi variabel seperti JWT_SECRET, DB_USER, DB_PASS bisa diakses via os.Getenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Mengecek JWT_SECRET terbaca atau tidak
	fmt.Println("JWT Secret:", os.Getenv("JWT_SECRET"))



	// Setup routes & jalankan server
	//Memanggil fungsi SetupRoutes() dari package routes, yang mengatur semua endpoint dan middleware (termasuk koneksi DB dan proteksi JWT/admin)
	r := routes.SetupRoutes()

	//Menampilkan pesan bahwa server aktif di port 8080
	//http.ListenAndServe akan menjalankan HTTP server di localhost:8080, dan memakai router r untuk menangani semua request
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
