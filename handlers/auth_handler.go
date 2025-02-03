package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "mobilerecharge/models"
)

func (h *Handler) Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var user models.User
    result := h.DB.Where("username = ? AND password = ?", loginData.Username, loginData.Password).First(&user)
    
    if result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Set cookie with proper settings
    c.SetCookie("logged_in", "true", 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "redirect": "/",
    })
}