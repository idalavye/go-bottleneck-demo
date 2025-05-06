// @title Go Bottlenecks API
// @description API for demonstrating Go performance bottlenecks
// @version 1.0.0
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idagdelen/go-bottlenecks/pkg/api"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	/*
		Create a trace file
	*/
	f, err := os.Create("macro-trace.out")
	if err != nil {
		log.Fatal(err)
	}
	if err := trace.Start(f); err != nil {
		log.Fatal(err)
	}

	/*
		Start pprof server
	*/
	go func() {
		log.Println("pprof started on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("Starting Go-Bottlenecks demonstration")
	fmt.Println("Welcome to Go performance bottlenecks demonstration")

	r := chi.NewRouter()

	/*
		Middleware
	*/
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	/*
		Routes
	*/
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	/*
		Serve Swagger UI using http-swagger package
	*/
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	/*
		Serve swagger.json directly from pkg/docs
	*/
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/docs/swagger.json")
	})

	/*
		API routes
	*/
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", api.HandleHealthCheck)
		r.Get("/search", api.HandleSearch)
	})

	/*
		Create a new HTTP server
	*/
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	/*
		Graceful shutdown
	*/
	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)

		trace.Stop()
		f.Close()
		close(idleConnsClosed)
	}()

	/*
		Start the server
	*/
	log.Println("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	/*
		Wait for the server to shut down
	*/
	<-idleConnsClosed
}
