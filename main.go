package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "gorm.io/gorm/logger"
    "mobilerecharge/models"
    "mobilerecharge/handlers"
    "mobilerecharge/config"
    "mobilerecharge/services"
    "time"
)

// Add this middleware function
func authRequired(c *gin.Context) {
    loggedIn, _ := c.Cookie("logged_in")
    if loggedIn != "true" {
        c.Redirect(http.StatusFound, "/login")
        c.Abort()
        return
    }
    c.Next()
}

func main() {
    // Set Gin to release mode before creating the engine
    gin.SetMode(gin.ReleaseMode)

    // Initialize database
    dsn := config.GetDBConfig()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        panic("failed to connect database: " + err.Error())
    }
    
    // Migrate the schema
    db.AutoMigrate(&models.Sim{}, &models.User{})  // Combine migrations
    
    // Create default admin user if not exists
    var adminUser models.User
    if db.Where("username = ?", "admin69").First(&adminUser).Error != nil {
        db.Create(&models.User{
            Username: "admin69",
            Password: "696969",
        })
    }

    // Initialize handler
    h := handlers.NewHandler(db)

    // Initialize notification service
    notificationService := services.NewNotificationService(db)

    // Start notification checker in a goroutine
    go func() {
        for {
            if err := notificationService.CheckAndSendNotifications(); err != nil {
                fmt.Printf("Error sending notifications: %v\n", err)
            }
            time.Sleep(12 * time.Hour)
        }
    }()

    // Initialize router after setting release mode
    r := gin.Default()

    // Add health check endpoint
    r.GET("/health", func(c *gin.Context) {
        // Check database connection
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "Database connection error"})
            return
        }
        
        // Ping database
        if err := sqlDB.Ping(); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "Database ping failed"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    // Serve static files
    r.Static("/static", "./static")
    r.LoadHTMLGlob("static/*.html")  // Update this line to load all HTML files
    
    // Public routes
    r.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)  // Use http.StatusOK constant
    })
    r.POST("/api/login", h.Login)

    // Protected routes
    protected := r.Group("/")
    protected.Use(authRequired)
    {
        protected.GET("/", func(c *gin.Context) {
            c.HTML(http.StatusOK, "index.html", nil)  // Use http.StatusOK constant
        })
        
        api := protected.Group("/api")
        {
            api.POST("/sims", h.AddSim)
            api.GET("/sims", h.GetAllSims)
            api.PUT("/sims/:id", h.UpdateSimRechargeDate)
        }
    }

    // Use config port directly
    port := config.GetPort()
    log.Printf("Server starting on port %s", port)
    r.Run("0.0.0.0:" + port)
}