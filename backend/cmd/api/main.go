package main

import (
	"fmt"
	"log"

	"github.com/Snehaaa-Kri/job-apply-automator/config"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/api"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/api/handlers"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/db"
	repoPostgres "github.com/Snehaaa-Kri/job-apply-automator/internal/repository/postgres"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/service"
	"github.com/Snehaaa-Kri/job-apply-automator/pkg/logger"
)

func main() {
	// 1. Load config from .env.local
	cfg, err := config.Load(".env.local")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// 2. Init logger
	log_, err := logger.New(cfg.App.Env)
	if err != nil {
		log.Fatalf("logger error: %v", err)
	}
	defer log_.Sync() // flush logs on exit

	log_.Info("starting job-apply-automator API")

	// 3. Connect to PostgreSQL
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		log_.Fatal("database connection failed", logger.Field("error", err))
	}
	log_.Info("connected to postgres")

	// 4. Wire dependencies (manual dependency injection)
	// Node analogy: const userRepo = new UserRepository(db)
	//               const authService = new AuthService(userRepo, config, logger)
	//               const authHandler = new AuthController(authService)
	userRepo    := repoPostgres.NewUserRepo(database)
	authSvc     := service.NewAuthService(userRepo, cfg, log_)
	authHandler := handlers.NewAuthHandler(authSvc)

	// 5. Set up router
	router := api.NewRouter(cfg, authHandler)

	// 6. Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log_.Info("server listening", logger.Field("addr", addr))

	if err := router.Engine().Run(addr); err != nil {
		log_.Fatal("server failed to start", logger.Field("error", err))
	}
}
