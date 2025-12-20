package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"stocky/internal/service"
)

type PortfolioHandler struct {
	portfolioService *service.PortfolioService
}

func NewPortfolioHandler(portfolioService *service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: portfolioService}
}

func (h *PortfolioHandler) GetTodayStocks(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	rewards, err := h.portfolioService.GetTodayRewards(c.Request.Context(), userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get today's rewards")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *PortfolioHandler) GetHistoricalINR(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	values, err := h.portfolioService.GetHistoricalINR(c.Request.Context(), userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get historical INR values")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, values)
}

func (h *PortfolioHandler) GetStats(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	stats, err := h.portfolioService.GetStats(c.Request.Context(), userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	portfolio, err := h.portfolioService.GetPortfolio(c.Request.Context(), userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get portfolio")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

