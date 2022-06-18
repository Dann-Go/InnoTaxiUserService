package handler

import (
	"net/http"

	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain/responses"
	"github.com/gin-gonic/gin"
)

// SignUp godoc
// @Summary      SignUp
// @Description  create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body domain.User true "User info"
// @Success      200  {object}  responses.ServerGoodResponse
// @Failure      400  {object}  responses.ServerError
// @Failure      404  {object}  responses.ServerError
// @Failure      500  {object}  responses.ServerError
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	json := domain.User{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewServerError(err.Error()))
		return
	}
	simplePassword := json.PasswordHash
	user, err := h.userService.CreateUser(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerError(err.Error()))
		return
	}

	token, err := h.authorizationService.GenerateToken(json.Phone, simplePassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewServerGoodResponse(map[string]interface{}{
		"user":  user,
		"token": token,
	}))
}

type signInInput struct {
	Phone        string `db:"phone" json:"phone" binding:"required"`
	PasswordHash string `db:"password_hash" json:"passwordHash" binding:"required"`
}

// SignIn godoc
// @Summary      SignIn
// @Description  Check if user exists and generate token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body signInInput true "User signIN info"
// @Success      200  {object}  responses.ServerGoodResponse
// @Failure      400  {object}  responses.ServerError
// @Failure      404  {object}  responses.ServerError
// @Failure      500  {object}  responses.ServerError
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	json := signInInput{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerError(err.Error()))
		return
	}
	token, err := h.authorizationService.GenerateToken(json.Phone, json.PasswordHash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerError(err.Error()))
		return
	}

	userFull, err := h.userService.GetUserByPhone(json.Phone)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewServerError(err.Error()))
		return
	}
	userResponse := domain.UserResponse{
		ID:     userFull.ID,
		Name:   userFull.Name,
		Phone:  userFull.Phone,
		Email:  userFull.Email,
		Rating: userFull.Rating,
	}

	c.JSON(http.StatusOK, responses.NewServerGoodResponse(map[string]interface{}{
		"user":  userResponse,
		"token": token,
	}))

}
