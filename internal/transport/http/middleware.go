package http

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (http *Http) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the JWT token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return echo.NewHTTPError(403, "Missing JWT token")
		}

		// Parse the JWT token and extract the ID
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Provide the JWT secret key for token verification
			return []byte("accesstokensecret"), nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(403, "Invalid JWT token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(403, "Invalid JWT claims")
		}

		userID, ok := claims["id"].(string)
		if !ok {
			return echo.NewHTTPError(403, "Invalid user ID in JWT claims")
		}

		// Set the ID as a request context value
		c.Set("UserID", userID)

		// Call the next handler
		return next(c)
	}
}
