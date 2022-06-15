package handler

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	json := domain.User{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewServerBadRequestError(err.Error()))
		return
	}
	simplePassword := json.PasswordHash
	user, err := h.services.CreateUser(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerInternalError(err.Error()))
		return
	}

	token, err := h.services.GenerateToken(json.Phone, simplePassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerInternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

type signInInput struct {
	Phone        string `db:"phone" json:"phone" binding:"required"`
	PasswordHash string `db:"password_hash" json:"passwordHash" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	json := signInInput{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerInternalError(err.Error()))
		return
	}
	token, err := h.services.GenerateToken(json.Phone, json.PasswordHash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerInternalError(err.Error()))
		return
	}

	user, err := h.services.GetUserByPhone(json.Phone)
	user.PasswordHash = ""

	c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})

}
