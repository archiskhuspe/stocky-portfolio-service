package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"stocky/internal/config"
	"stocky/internal/database"
	"stocky/internal/handler"
	"stocky/internal/middleware"
	"stocky/internal/repository"
	"stocky/internal/scheduler"
	"stocky/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	gin.SetMode(cfg.Server.GinMode)

	db, err := database.NewPostgres(cfg.Database.DSN())
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	rewardRepo := repository.NewRewardRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	priceRepo := repository.NewStockPriceRepository(db)

	priceService := service.NewPriceService(priceRepo)
	rewardService := service.NewRewardService(rewardRepo, ledgerRepo, userRepo, priceRepo, db)
	portfolioService := service.NewPortfolioService(rewardRepo, priceRepo)

	rewardHandler := handler.NewRewardHandler(rewardService)
	portfolioHandler := handler.NewPortfolioHandler(portfolioService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	api := router.Group("/api/v1")
	{
		api.POST("/reward", rewardHandler.CreateReward)
		api.GET("/today-stocks/:userId", portfolioHandler.GetTodayStocks)
		api.GET("/historical-inr/:userId", portfolioHandler.GetHistoricalINR)
		api.GET("/stats/:userId", portfolioHandler.GetStats)
		api.GET("/portfolio/:userId", portfolioHandler.GetPortfolio)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priceFetcher := scheduler.NewPriceFetcher(priceService, cfg.PriceService.FetchInterval)
	go priceFetcher.Start(ctx)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		logrus.WithField("port", cfg.Server.Port).Info("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server forced to shutdown")
	}

	logrus.Info("Server exited")
}

