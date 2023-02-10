package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/configs"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
)

type token struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
} // @Name Token

type resetPassword struct {
	ResetToken  string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"gte=8"`
} // @Name ResetPassword

// @Tags Account
// @Summary Current
// @Description Read current user
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/account/current [get]
// @Security ApiKeyAuth
func Current(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	return c.Status(http.StatusOK).JSON(currentUser)
}

// @Tags Account
// @Summary Login
// @Description Login
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} token
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/account/login [post]
func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	dbUser, err := crud.AuthenticateUser(username, password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Error{Detail: "Incorrect email or password"})
	}

	tokenString, err := utils.GetJwt(dbUser.ID.Hex(), configs.Env.AccessTokenExpirationMinutes)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(token{AccessToken: tokenString, TokenType: "Bearer"})
}

// @Tags Account
// @Summary Recover Password
// @Description Recover password
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/account/recover-password/{email} [post]
func RecoverPassword(c *fiber.Ctx) error {
	email := c.Params("email")

	_, err := crud.FindOneUserByEmail(email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "User not found"})
	}

	utils.SendResetPasswordEmail(email)

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: "Password recovery email sent"})
}

// @Tags Account
// @Summary Reset Password
// @Description Reset password
// @Accept json
// @Produce json
// @Param body body resetPassword true "Body"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/account/reset-password [post]
func ResetPassword(c *fiber.Ctx) error {
	body := resetPassword{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Error{Detail: err.Error()})
	}

	if err := utils.NewValidator().Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	claims, err := utils.GetJwtClaims(body.ResetToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: "Failed to validate credentials"})
	}

	dbUser, err := crud.FindOneUserByEmail(claims.Subject)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "User not found"})
	}

	if !dbUser.IsActive {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: "The user is inactive"})
	}

	dbUser.Password = body.NewPassword

	dbUser, err = crud.UpdateUser(dbUser.ID.Hex(), dbUser)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: "Password updated successfully"})
}
