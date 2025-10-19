package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harlesbayu/fleet-backend/internal/usecase"
)

type VehicleHandler struct {
	usecase *usecase.VehicleUsecase
}

func NewVehicleHandler(u *usecase.VehicleUsecase) *VehicleHandler {
	return &VehicleHandler{usecase: u}
}

func (h *VehicleHandler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	loc, err := h.usecase.GetLastLocation(vehicleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, loc)
}

func (s *VehicleHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	from, _ := strconv.ParseInt(c.Query("start"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("end"), 10, 64)

	data, err := s.usecase.GetHistory(vehicleID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
