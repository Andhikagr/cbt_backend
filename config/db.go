package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//Variabel global DB akan menyimpan koneksi database yang bisa dipakai di seluruh aplikasi (misalnya di handlers atau models)
var DB *sql.DB

// Variabel global JWTSecret menyimpan token rahasia untuk keperluan JWT. Disarankan tidak langsung mengisi dari os.Getenv(), tapi tunggu setelah .env berhasil dimuat dulu di InitDB()
var JWTSecret string // Jangan langsung isi dari os.Getenv()


//Membaca file .env, mengisi JWTSecret, membuat koneksi ke MySQL dan menyimpannya ke DB
func InitDB() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env:", err)
	}

	// Ambil JWT secret dari .env setelah .env berhasil dimuat
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET belum diset di file .env")
	}

	// Ambil konfigurasi DB dari .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Buat DSN string koneksi
	//Menyusun DSN (Data Source Name) untuk koneksi MySQL lokal di 127.0.0.1:3306. Tambahan ?parseTime=true agar field DATETIME bisa dibaca sebagai time.Time di Go
	dsn := dbUser + ":" + dbPass + "@tcp(127.0.0.1:3306)/" + dbName + "?parseTime=true"

	// Inisialisasi koneksi database, membuka koneksi ke database dengan driver mysql dan DSN yang sudah disusun
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
