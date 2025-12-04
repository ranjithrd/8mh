package middleware

import (
	"backend/src/repos"
	"net/http"

	"github.com/labstack/echo/v4"
)

var sessionRepo = repos.SessionRepo{}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		session, err := sessionRepo.FindBySessionID(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		}

		userWithSession := &repos.UserWithSession{
			ID:             session.User.ID,
			PhoneNumber:    session.User.PhoneNumber,
			Name:           session.User.Name,
			Email:          session.User.Email,
			SavingsBalance: session.User.SavingsBalance,
			SharesBalance:  session.User.SharesBalance,
			IsActive:       session.User.IsActive,
		}

		c.Set("user", userWithSession)
		return next(c)
	}
}
