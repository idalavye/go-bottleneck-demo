// @title Go Bottlenecks API
// @description API for demonstrating Go performance bottlenecks
// @version 1.0.0
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idagdelen/go-bottlenecks/pkg/api"
	_ "github.com/idagdelen/go-bottlenecks/pkg/docs" // Swagger docs içe aktarılıyor
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	log.Println("Starting Go-Bottlenecks demonstration")
	fmt.Println("Welcome to Go performance bottlenecks demonstration")

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// Serve Swagger UI using http-swagger package
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), // Swagger dokümantasyonunun URL'i
	))

	// Serve swagger.json directly from pkg/docs
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/docs/swagger.json")
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", api.HandleHealthCheck)
		r.Get("/sequential", api.HandleSequential)
		r.Get("/concurrent", api.HandleConcurrent)
		r.Get("/pool", api.HandlePool)
		r.Get("/leak", api.HandleLeak)
		r.Get("/search", api.HandleSearch)
	})

	// Start server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
