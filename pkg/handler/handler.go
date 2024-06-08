package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"hb-notification/pkg/service"
)

type Handler struct {
	services *service.Service
	Cron     *cron.Cron
	JobID    cron.EntryID
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	notification := router.Group("/", h.userIdentity)
	{
		notification.GET("/users", h.getUsers)
		notification.GET("/subscriptions", h.getSubscriptions)
		notification.POST("/subscribe", h.subscribe)
		notification.DELETE("/unsubscribe", h.unsubscribe)
	}

	router.POST("/update-cron", h.updateCron)

	return router
}
