package validator

import (
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewValidator() *CustomValidator {
	cv := &CustomValidator{
		Validator: validator.New(),
	}
	cv.init()
	return cv
}

func (cv *CustomValidator) init() {
	cv.Validator.RegisterValidation("password", password)
	cv.Validator.RegisterValidation("phoneNumber", phoneNumber)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func password(fl validator.FieldLevel) bool {
	var (
		hasNumber      = false
		hasSpecialChar = false
		hasLetter      = false
		hasSuitableLen = false
	)

	password := fl.Field().String()

	if utf8.RuneCountInString(password) <= 64 && utf8.RuneCountInString(password) >= 6 {
		hasSuitableLen = true
	}

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecialChar = true
		case unicode.IsLetter(c) || c == ' ':
			hasLetter = true
		default:
			return false
		}
	}

	return hasNumber && hasSpecialChar && hasLetter && hasSuitableLen
}

func phoneNumber(lf validator.FieldLevel) bool {
	phoneNumber := lf.Field().String()
	if len(phoneNumber) < 9 || len(phoneNumber) > 14 {
		return false
	}

	return phoneNumber[:3] == "+62"
}
