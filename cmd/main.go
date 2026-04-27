package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ppablomunoz/go-shorten/internal/handler"
	_ "modernc.org/sqlite"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		log.Fatalf("Error opening DB. Path: %s\n", dbURL)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot connect to db. Path: %s\n", dbURL)
	}

	// Initialize database schema
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Printf("Warning: could not read schema.sql: %v", err)
	} else {
		_, err = db.Exec(string(schema))
		if err != nil {
			log.Fatalf("Error initializing database schema: %v", err)
		}
		log.Println("Database schema initialized successfully")
	}

	handler := handler.NewHandler(db)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Load templates
	r.LoadHTMLGlob("web/templates/*")
	// Static files
	r.Static("/static", "./web/static")

	// ====== Frontend pages ======
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/:code", handler.EnterURL)

	api := r.Group("/api")

	api.GET("/health", func(ctx *gin.Context) {
		if err := db.Ping(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "database connection lost"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ====== Backend endpoints ======
	api.POST("/url", handler.NewURL)
	api.GET("/url", handler.GetURLs)
	api.PUT("/url/:code", handler.UpdateURL)
	api.DELETE("/url/:code", handler.DeleteURL)

	// Configure the HTTP server
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server in a separate goroutine so it doesn't block the shutdown signal listening below
	go func() {
		log.Printf("Running on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ============================
	// ====== Greateful Exit ======
	// ============================

	// Create a channel to listen for OS interrupt signals (like Ctrl+C or Docker stop)
	quit := make(chan os.Signal, 1)
	// SIGINT: Ctrl+C, SIGTERM: Generic termination signal (used by Docker/K8s)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block execution here until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for the shutdown process (5 seconds)
	// This ensures the server doesn't hang forever if a request is stuck
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop the server:
	// 1. Stop accepting new connections
	// 2. Wait for active requests to finish within the 5s timeout
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
