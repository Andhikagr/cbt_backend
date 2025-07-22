package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

var JWTSecret string 

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env:", err)
	}


	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET belum diset di file .env")
	}

	// Ambil konfigurasi DB dari .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	
	dsn := dbUser + ":" + dbPass + "@tcp(127.0.0.1:3306)/" + dbName + "?parseTime=true"


	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Cek koneksi
	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal ping ke database:", err)
	}

	log.Println("Berhasil koneksi ke database MySQL")
}
