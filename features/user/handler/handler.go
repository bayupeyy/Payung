package handler

import (
	"21-api/features/user"
	"21-api/helper"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service user.Service
}

func NewUserHandler(s user.Service) user.Controller {
	return &controller{
		service: s,
	}
}

// Register User
func (ct *controller) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input user.Register
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupported") {
				return c.JSON(http.StatusUnsupportedMediaType, helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.ErrorUserInputFormat, nil))
			}
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.ErrorUserInput, nil))
		}

		err = ct.service.Register(input)
		if err != nil {
			return c.JSON(helper.ErrorCode(err), helper.ResponseFormat(helper.ErrorCode(err), err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Registered Successfully", nil))
	}
}

func (ct *controller) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupported") {
				return c.JSON(http.StatusUnsupportedMediaType, helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.ErrorUserInputFormat, nil))
			}
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.ErrorUserInput, nil))
		}

		var processData user.User
		processData.Email = input.Email
		processData.Password = input.Password

		usertoken, err := ct.service.Login(processData)
		if err != nil {
			return c.JSON(helper.ErrorCode(err), helper.ResponseFormat(helper.ErrorCode(err), err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Login Successfully", usertoken))
	}
}

func (ct *controller) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.ErrorUserInput, nil))
		}

		profile, err := ct.service.Profile(token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ErrorGeneralServer, nil))
		}

		var profileResponse ProfileResponse
		profileResponse.Name = profile.Name
		profileResponse.Email = profile.Email
		profileResponse.Hp = profile.Hp

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully Get MyProfile", profileResponse))
	}
}

// Update User
func (ct *controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input user.User
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType, helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.ErrorUserInputFormat, nil))
			}
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.ErrorUserInput, nil))
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.ErrorUserInput, nil))
		}

		err = ct.service.Update(token, input)
		if err != nil {
			return c.JSON(helper.ErrorCode(err), helper.ResponseFormat(helper.ErrorCode(err), err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully Updated", nil))
	}
}

// Delete User
func (ct *controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.ErrorUserInput, nil))
		}

		err := ct.service.Delete(token)
		if err != nil {
			return c.JSON(helper.ErrorCode(err), helper.ResponseFormat(helper.ErrorCode(err), err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully Deleted User", nil))
	}
}
