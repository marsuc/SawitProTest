package handler

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/marsuc/SawitProTest/pkg/jwtrsa"
)

func GetTokenClaims(ctx echo.Context) (token *jwt.Token, err error) {
	bearerToken := ctx.Request().Header.Get("Authorization")
	if bearerToken == "" {
		return nil, errors.New("no bearer token")
	}

	bearerToken = bearerToken[len("Bearer "):]
	return jwtrsa.ValidateJWT(publicKey, bearerToken)
}
