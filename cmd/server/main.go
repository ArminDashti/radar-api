package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ArminDashti/radar-api/internal/config"
	"github.com/ArminDashti/radar-api/internal/db"
	"github.com/ArminDashti/radar-api/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := db.RunMigration(ctx, pool, "migrations/001_init.sql"); err != nil {
		log.Fatal(err)
	}
	if err := db.Seed(ctx, pool); err != nil {
		log.Fatal(err)
	}

	server := &handlers.Server{Pool: pool, JWTSecret: cfg.JWTSecret}
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	api.POST("/auth/login", server.Login)
	web := api.Group("")
	web.Use(server.WebAuth())
	web.GET("/probes", server.ListProbes)
	web.GET("/endpoints", server.ListEndpoints)
	web.POST("/endpoints", server.CreateEndpoint)
	web.GET("/grid/endpoints", server.EndpointGrid)
	web.GET("/grid/probes", server.ProbeGrid)
	agent := api.Group("/agent")
	agent.Use(server.AgentAuth())
	agent.GET("/targets", server.AgentTargets)
	agent.POST("/samples", server.AgentSamples)

	log.Printf("radar API listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
