package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/marsuc/SawitProTest/generated"
	"github.com/marsuc/SawitProTest/pkg/password"
	"github.com/marsuc/SawitProTest/repository"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostRegister(ctx echo.Context) error {
	errorResponse := generated.ErrorResponse{}
	var req generated.RegisterRequest
	err := ctx.Bind(&req)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	user := repository.User{}.FromRegisterRequest(req)
	err = ctx.Validate(user)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	user.Password, err = password.HashPassword(user.Password)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	existingUser, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	if existingUser.PhoneNumber != "" {
		errorResponse.Message = "phone number already exists"
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	userId, err := s.Repository.CreateUser(ctx.Request().Context(), user)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	resp := generated.RegisterResponse{
		Id: userId,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	errorResponse := generated.ErrorResponse{}
	var req generated.LoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	user, valid, err := s.validateLogin(ctx.Request().Context(), req.PhoneNumber, req.Password)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	if !valid {
		errorResponse.Message = invalidUsernameOrPassword
		return ctx.JSON(http.StatusUnauthorized, errorResponse)
	}

	accessToken, err := s.createAccessToken(ctx.Request().Context(), user)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	err = s.updateLoginCount(ctx.Request().Context(), user)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	resp := generated.LoginResponse{
		UserId:      user.Id,
		AccessToken: accessToken,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetProfile(ctx echo.Context) error {
	token, err := GetTokenClaims(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp := generated.GetProfileResponse{
		FullName:    token.Claims.(jwt.MapClaims)["full_name"].(string),
		PhoneNumber: token.Claims.(jwt.MapClaims)["phone_number"].(string),
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	token, err := GetTokenClaims(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	errorResponse := generated.ErrorResponse{}
	var req generated.UpdateProfileRequest
	err = ctx.Bind(&req)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	user, err := s.Repository.GetUserById(ctx.Request().Context(), int(token.Claims.(jwt.MapClaims)["id"].(float64)))
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	existingUser, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	if existingUser.Id != user.Id && existingUser.Id != 0 {
		errorResponse.Message = "phone number already exists"
		return ctx.JSON(http.StatusConflict, errorResponse)
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	user.PhoneNumber = req.PhoneNumber

	err = ctx.Validate(user)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	err = s.Repository.UpdateUser(ctx.Request().Context(), user)
	if err != nil {
		errorResponse.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errorResponse)
	}

	resp := generated.UpdateProfileJSONRequestBody{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}
