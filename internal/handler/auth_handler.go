package handler

import (
	"net/http"

	"github.com/C0deNeo/goSessionStore/internal/usercase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Usecase *usercase.AuthUseCase
}

func NewAuthHandler(u *usercase.AuthUseCase) *AuthHandler {
	return &AuthHandler{u}
}

func (h *AuthHandler) Signup(c echo.Context) error {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Req

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	err := h.Usecase.SignUp(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "signup successful"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Req

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})

	}
	token, err := h.Usecase.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing token"})
	}

	err := h.Usecase.Logout(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "logged out successfully!"})
}
