package middleware

import (
	"net/http"

	"github.com/C0deNeo/goSessionStore/internal/pkg/jwt"
	"github.com/C0deNeo/goSessionStore/internal/repository"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(sessionRepo *repository.RedisSessionRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			//get the token from the request header
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
			}

			//parse the token
			_, err := jwt.ParseToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			//check for the token validity
			valid, err := sessionRepo.IsTokenValid(c.Request().Context(), token)
			if err != nil || !valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "session expired"})
			}
			return next(c)
		}
	}
}
