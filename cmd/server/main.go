package main

import (
	"log"
	"net/http"

	"example.com/mod/internal/database"
	"example.com/mod/internal/handler"
	"example.com/mod/internal/middleware"
	"example.com/mod/internal/postgres"
	"example.com/mod/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment")
	}

	err = database.Connection()
	if err != nil {
		log.Fatal(err)
	}

	defer database.CloseConnection()

	pool := database.GetDB()
	jobRepo := postgres.NewPostgresJobRepository(pool)
	userRepo := postgres.NewPostgresUserRepository(pool)
	authRepo := postgres.NewPostgresAuthRepository(pool)
	authService := service.NewAuthService(userRepo, authRepo)
	authHandler := handler.NewAuthHandler(authService)
	JobService := service.NewJobService(jobRepo)
	JobHandler := handler.NewJobHandler(JobService)

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.Get("/jobs", JobHandler.GetJobs)
	r.Get("/jobs/{id}", JobHandler.GetJob)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)
	r.Get("/health", handler.HealthCheck)
	r.Get("/auth/me", authHandler.Me)
	r.Post("/auth/logout", authHandler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RequireRole("employer"))
		r.Post("/jobs", JobHandler.PostJob)
		r.Put("/jobs/{id}", JobHandler.PutJob)
		r.Delete("/jobs/{id}", JobHandler.DeleteJob)
	})

	log.Fatal(http.ListenAndServe(":8090", r))
}
