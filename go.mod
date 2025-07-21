module cbt-backend

go 1.24.5


//Ini adalah library MySQL driver untuk Go (go-sql-driver/mysql), versi 1.9.3. Digunakan untuk menghubungkan aplikasi Go ke database MySQL

require github.com/go-sql-driver/mysql v1.9.3

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.3 // Untuk membuat dan memverifikasi JWT token di fitur login
	github.com/gorilla/mux v1.8.1 // Router HTTP yang fleksibel dan powerful untuk menangani endpoint AP
	
	//Untuk membaca file .env agar variabel environment bisa digunakan di aplikasi Go
	github.com/joho/godotenv v1.5.1 

	//Kumpulan algoritma enkripsi, termasuk bcrypt untuk hashing password
	golang.org/x/crypto v0.40.0
)
