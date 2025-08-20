package main

import (
    "log"

    "loan_service/internal/config"
    "loan_service/internal/domain"
    "loan_service/internal/handler"
    "loan_service/internal/repository"
    "loan_service/internal/service"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // Load configuration from environment variables
    cfg := config.Load()

    // Initialize database connection
    db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }
    // Enable UUID extension for Postgres so gorm can generate UUIDs
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
        log.Fatalf("failed to create uuid extension: %v", err)
    }
    // Auto migrate the schema. This will create tables if they do not
    // exist. In production you may want to manage migrations using
    // dedicated tools instead of gorm's automigrate.
    if err := db.AutoMigrate(
        &domain.Loan{},
        &domain.Approval{},
        &domain.Investor{},
        &domain.Investment{},
        &domain.Disbursement{},
    ); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    // Initialize repository, service and handlers
    repo := repository.NewLoanRepository(db)
    svc := service.NewLoanService(repo)
    loanHandler := handler.NewLoanHandler(svc)

    // Configure Gin router
    r := gin.Default()
    loanHandler.RegisterRoutes(r)

    // Start HTTP server
    addr := ":" + cfg.ServerPort
    log.Printf("starting server at %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("server error: %v", err)
    }
}