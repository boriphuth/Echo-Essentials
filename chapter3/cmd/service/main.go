package main

import (
	"github.com/PacktPublishing/Echo-Essentials/chapter3/handlers"
	"github.com/PacktPublishing/Echo-Essentials/chapter3/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// create a new echo instance
	e := echo.New()

	// Signing Key for our auth middleware
	var signingKey = []byte("superdupersecret!")
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(models.SigningContextKey, signingKey)
			return next(c)
		}
	})

	// in order to serve static assets
	e.Static("/static", "static")

	// reminder handler group
	reminderGroup := e.Group("/reminder")
	reminderGroup.Use(middleware.JWT(signingKey))
	reminderGroup.POST("", handlers.CreateReminder)

	// Route / to handler function
	e.GET("/health-check", handlers.HealthCheck)
	// Authentication routes
	e.POST("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)

	// V1 Routes
	v1 := e.Group("/v1")
	// V1 Authentication routes
	v1.POST("/login", handlers.Login)
	v1.POST("/logout", handlers.Logout)
	// V1 Reminder Routes
	v1Reminders := v1.Group("/reminder", middleware.JWT(signingKey))
	v1Reminders.POST("", handlers.CreateReminder)
	v1Reminders.GET("/:id", handlers.GetReminder)
	// /v1/reminder/:id

	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":8080"))
}
