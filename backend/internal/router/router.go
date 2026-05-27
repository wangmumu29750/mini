package router

import (
	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/handler"
	"mini-12306/backend/internal/middleware"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"mini-12306/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(middleware.TraceID())
	r.Use(gin.Logger())
	r.Use(middleware.Recovery())

	api := r.Group("/api/v1")
	{
		healthHandler := handler.NewHealthHandler(cfg, db)
		api.GET("/health", healthHandler.Check)

		userRepo := repository.NewUserRepository(db)
		trainRepo := repository.NewTrainRepository(db)
		orderRepo := repository.NewOrderRepository(db)
		ticketRepo := repository.NewTicketRepository(db)
		settingRepo := repository.NewSystemSettingRepository(db)
		authService := service.NewAuthService(cfg, userRepo)
		trainService := service.NewTrainService(trainRepo)
		orderService := service.NewOrderService(cfg, orderRepo)
		ticketService := service.NewTicketService(ticketRepo)
		settingService := service.NewSystemSettingService(settingRepo)
		authHandler := handler.NewAuthHandler(authService)
		trainHandler := handler.NewTrainHandler(trainService)
		orderHandler := handler.NewOrderHandler(orderService)
		ticketHandler := handler.NewTicketHandler(ticketService)
		settingHandler := handler.NewSystemSettingHandler(settingService)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/me", middleware.AuthRequired(cfg.JWTSecret), authHandler.Me)
		}

		api.GET("/stations", trainHandler.Stations)
		api.GET("/trains/search", trainHandler.Search)

		orderGroup := api.Group("/orders", middleware.AuthRequired(cfg.JWTSecret))
		{
			orderGroup.POST("", orderHandler.Create)
			orderGroup.GET("", orderHandler.List)
			orderGroup.GET("/:id", orderHandler.Detail)
			orderGroup.POST("/:id/cancel", orderHandler.Cancel)
			orderGroup.POST("/:id/payments", orderHandler.Pay)
		}

		ticketGroup := api.Group("/tickets", middleware.AuthRequired(cfg.JWTSecret))
		{
			ticketGroup.GET("", ticketHandler.List)
			ticketGroup.GET("/:id", ticketHandler.Detail)
			ticketGroup.GET("/:id/change-options", ticketHandler.ChangeOptions)
			ticketGroup.POST("/:id/refund", ticketHandler.Refund)
			ticketGroup.POST("/:id/change", ticketHandler.Change)
		}

		clerkGroup := api.Group("/clerk", middleware.AuthRequired(cfg.JWTSecret), middleware.RoleRequired(string(model.UserRoleClerk), string(model.UserRoleAdmin)))
		{
			clerkGroup.POST("/orders", orderHandler.ClerkCreate)
		}

		adminGroup := api.Group("/admin", middleware.AuthRequired(cfg.JWTSecret), middleware.RoleRequired(string(model.UserRoleAdmin)))
		{
			adminGroup.GET("/settings", settingHandler.List)
			adminGroup.PUT("/settings", settingHandler.Update)
		}
	}

	return r
}
