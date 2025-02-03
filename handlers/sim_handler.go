package handlers

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "mobilerecharge/models"
)

type Handler struct {
    DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{DB: db}
}

func (h *Handler) AddSim(c *gin.Context) {
    var sim models.Sim
    if err := c.ShouldBindJSON(&sim); err != nil {
        fmt.Printf("Error binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("Received SIM data: %+v\n", sim)

    // Truncate time component to get only date
    sim.LastRechargeDate = sim.LastRechargeDate.Truncate(24 * time.Hour)
    
    // Calculate validities (only dates)
    sim.RechargeValidity = sim.LastRechargeDate.AddDate(0, 0, 28).Truncate(24 * time.Hour)
    sim.IncomingValidity = sim.RechargeValidity.AddDate(0, 0, 7).Truncate(24 * time.Hour)
    sim.SimExpiry = sim.IncomingValidity.AddDate(0, 0, 85).Truncate(24 * time.Hour)

    // Save to database
    result := h.DB.Create(&sim)
    if result.Error != nil {
        fmt.Printf("Error saving to database: %v\n", result.Error)
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    fmt.Printf("Successfully saved SIM with ID: %d\n", sim.ID)
    c.JSON(http.StatusOK, sim)
}

func (h *Handler) GetAllSims(c *gin.Context) {
    var sims []models.Sim
    result := h.DB.Order("last_recharge_date desc").Find(&sims)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, sims)
}

// Add this new handler function
func (h *Handler) UpdateSimRechargeDate(c *gin.Context) {
    type UpdateRequest struct {
        LastRechargeDate string `json:"lastRechargeDate" binding:"required"`
    }

    simID := c.Param("id")
    var updateData UpdateRequest
    
    if err := c.ShouldBindJSON(&updateData); err != nil {
        fmt.Printf("Binding error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parse the date string
    lastRechargeDate, err := time.Parse("2006-01-02", updateData.LastRechargeDate)
    if err != nil {
        fmt.Printf("Date parsing error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
        return
    }

    var sim models.Sim
    if err := h.DB.First(&sim, simID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "SIM not found"})
        return
    }

    // Update dates
    sim.LastRechargeDate = lastRechargeDate
    sim.RechargeValidity = lastRechargeDate.AddDate(0, 0, 28)
    sim.IncomingValidity = sim.RechargeValidity.AddDate(0, 0, 7)
    sim.SimExpiry = sim.IncomingValidity.AddDate(0, 0, 85)

    if err := h.DB.Save(&sim).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, sim)
}