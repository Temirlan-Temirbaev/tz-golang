package controllers

import (
	"github.com/Temirlan-Temirbaev/tz-golang/config"
	"github.com/Temirlan-Temirbaev/tz-golang/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

var oauthDataMap = make(map[string]struct {
	RequestSecret string
	TelegramID    string
})

func InitUserRoutes(e *echo.Echo) {
	e.GET("/users/login", func(c echo.Context) error {
		telegramId := c.QueryParam("telegram_id")
		requestToken, requestSecret, err := config.AuthConfig.RequestToken()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch request token"})
		}
		oauthDataMap[requestToken] = struct {
			RequestSecret string
			TelegramID    string
		}{
			RequestSecret: requestSecret,
			TelegramID:    telegramId,
		}
		authorizationURL, err := config.AuthConfig.AuthorizationURL(requestToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate authorization URL"})
		}
		return c.Redirect(http.StatusFound, authorizationURL.String())
	})
	e.GET("/users/callback", func(c echo.Context) error {
		requestToken := c.QueryParam("oauth_token")
		verifier := c.QueryParam("oauth_verifier")
		data, exists := oauthDataMap[requestToken]
		if !exists {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or expired request token"})
		}
		requestSecret := data.RequestSecret
		telegramId := data.TelegramID
		if !exists {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or expired request token"})
		}
		accessToken, accessSecret, err := config.AuthConfig.AccessToken(requestToken, requestSecret, verifier)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch access token"})
		}
		user := new(models.User)
		user.TelegramID = telegramId
		user.AccessSecret = accessSecret
		user.AccessToken = accessToken
		if err := config.DB.Create(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"AccessToken":  accessToken,
			"AccessSecret": accessSecret,
		})
	})
	e.GET("/users/id/:id", func(c echo.Context) error {
		id := c.Param("id")
		user := new(models.User)

		if err := config.DB.Preload("Tasks").Where("telegram_id = ?", id).First(user).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, user)
	})
}
