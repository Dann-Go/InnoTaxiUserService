package handler

import (
	_ "github.com/Dann-Go/InnoTaxiUserService/docs"
	"github.com/Dann-Go/InnoTaxiUserService/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	userService         *service.UserService
	authariztionService *service.AuthorizationService
}

func NewHandler(userService *service.UserService, authorizationService *service.AuthorizationService) *Handler {
	return &Handler{
		userService:         userService,
		authariztionService: authorizationService,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

	return router
}
