package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"stocky/internal/service"
)

type RewardHandler struct {
	rewardService *service.RewardService
}

func NewRewardHandler(rewardService *service.RewardService) *RewardHandler {
	return &RewardHandler{rewardService: rewardService}
}

func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req service.RewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.rewardService.ProcessReward(c.Request.Context(), req); err != nil {
		logrus.WithError(err).Error("Failed to process reward")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Reward processed successfully",
		"event_id": req.EventID,
	})
}

