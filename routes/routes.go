package routes

import (
	"cbt-backend/config"
	"cbt-backend/handlers"
	"cbt-backend/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Ganti * dengan domain asal untuk produksi
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}





func SetupRoutes() http.Handler {
	config.InitDB()

	r := mux.NewRouter()
	r.Use(corsMiddleware) 

	// ====== Public Routes ======
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to CBT API"))
	}).Methods("GET")

	// r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST", "OPTIONS")


	// Public GET Endpoints (tanpa token)
	r.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", handlers.GetQuestionByID).Methods("GET")
	r.HandleFunc("/exams", handlers.GetExams).Methods("GET")
	r.HandleFunc("/exams/{id}", handlers.GetExamByID).Methods("GET")
	r.HandleFunc("/subjects", handlers.GetAllSubjects).Methods("GET")
	r.HandleFunc("/subjects/{id}", handlers.GetSubjectByID).Methods("GET")
	r.HandleFunc("/grades", handlers.GetGrades).Methods("GET")
	r.HandleFunc("/grades/{id}", handlers.GetGradeByID).Methods("GET")
	r.HandleFunc("/answers", handlers.GetAnswers).Methods("GET")
	r.HandleFunc("/answers/{id}", handlers.GetAnswerByID).Methods("GET")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserByID).Methods("GET")

	// ====== JWT Protected Routes (Login Required) ======
	authRoutes := r.PathPrefix("/").Subrouter()
	authRoutes.Use(middlewares.JWTMiddleware)

	authRoutes.HandleFunc("/questions", handlers.CreateQuestion).Methods("POST")
	authRoutes.HandleFunc("/questions/{id}", handlers.UpdateQuestionByID).Methods("PUT")
	authRoutes.HandleFunc("/questions/{id}", handlers.DeleteQuestionByID).Methods("DELETE")

	authRoutes.HandleFunc("/exams", handlers.CreateExam).Methods("POST")
	authRoutes.HandleFunc("/exams/{id}", handlers.UpdateExam).Methods("PUT")
	authRoutes.HandleFunc("/exams/{id}", handlers.DeleteExam).Methods("DELETE")

	authRoutes.HandleFunc("/subjects", handlers.CreateSubject).Methods("POST")
	authRoutes.HandleFunc("/subjects/{id}", handlers.DeleteSubject).Methods("DELETE")

	authRoutes.HandleFunc("/grades", handlers.CreateGrade).Methods("POST")
	authRoutes.HandleFunc("/grades/{id}", handlers.DeleteGrade).Methods("DELETE")

	// ====== Admin Only Routes ======
	adminRoutes := authRoutes.PathPrefix("/").Subrouter()
	adminRoutes.Use(middlewares.AdminOnly)

	adminRoutes.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	adminRoutes.HandleFunc("/register", handlers.RegisterUser).Methods("POST")

	return r
}
