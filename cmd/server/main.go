package main

import (
	"log"
	"net/http"

	"example.com/mod/internal/auth"
	"example.com/mod/internal/database"
	"example.com/mod/internal/handlers"
	"example.com/mod/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment")
	}

	error := database.Connection()
	if error != nil {
		log.Fatal(error)
	}

	defer database.CloseConnection()

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// r.Get("/bookmarks", handlers.GetBooks)
	// r.Post("/bookmarks", handlers.PostBooks)

	// r.Get("/bookmarks/{id}", handlers.GetOneBookmark)
	// r.Put("/bookmarks/{id}", handlers.UpdateBookmark)
	// r.Delete("/bookmarks/{id}", handlers.DeleteBookmark)

	// This for different handler
	r.Get("/jobs", handlers.GetJobs)
	r.Get("/jobs/{id}", handlers.GetJob)
	r.Post("/auth/register", auth.Register)
	r.Post("/auth/login", auth.Login)
	r.Post("/auth/refresh", auth.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RequireRole("employer"))
		r.Post("/jobs", handlers.PostJob)
		r.Put("/jobs/{id}", handlers.PutJob)
		r.Delete("/jobs/{id}", handlers.DeleteJob)
	})

	http.ListenAndServe(":8090", r)
}
