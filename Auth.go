/*

LOGIN REQUEST:
http://127.0.0.1:3543/Login

{
	"Username": "Necdet",
	"Password": "12345alti"
}

---

REQUEST HEADER:
Authorization:Bearer <TOKEN>

*/

package main

import (
	"goapi_base/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AuthRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type AuthResponse struct {
	Token string `json:"Token"`
}

func Login(c echo.Context) error {
	u := &AuthRequest{}
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.Username != config.AUTH_USERNAME || u.Password != config.AUTH_PASSWORD {
		return echo.ErrUnauthorized
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "JWT_CLAIM"
	claims["exp"] = time.Now().Add(time.Hour * config.AUTH_TIMEOUT_HOUR).Unix()
	t, err := token.SignedString([]byte(config.AUTH_JWT_TOKEN))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
