package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/config"
	"github.com/wilfredohq/fiber-start/constants"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Tags Account
// @Summary Get Current Account
// @Description Get current account
// @Accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure default {object} models.Error
// @Router /api/v1/account/current [get]
// @Security ApiKeyAuth
func GetCurrentAccount(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	return c.Status(http.StatusOK).JSON(currentUser)
}

type TokenResponse struct {
	AccessToken string `json:"accessToken" validate:"required"`
	TokenType   string `json:"tokenType" validate:"required"`
} // @Name Token

// @Tags Account
// @Summary Login
// @Description Login
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} TokenResponse
// @Failure default {object} models.Error
// @Router /api/v1/account/login [post]
func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	userResponse, err := crud.AuthenticateUser(username, password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: constants.InvalidCredentials})
	}

	tokenString, err := utils.GetJwt(userResponse.ID.Hex(), config.Config.AccessTokenExpirationMinutes)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(TokenResponse{AccessToken: tokenString, TokenType: "Bearer"})
}

type RecoverAccountBody struct {
	Email string `json:"email" validate:"required,email"`
} // @Name RecoverAccount

// @Tags Account
// @Summary Recover Account
// @Description Recover account
// @Accept json
// @Produce json
// @Param body body RecoverAccountBody true "Body"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/account/recover [post]
func RecoverAccount(c *fiber.Ctx) error {
	body := RecoverAccountBody{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	if _, err := crud.FindOneUserByEmail(body.Email); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.UserNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	utils.SendResetPasswordEmail(body.Email)

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: constants.EmailSent})
}

type ResetPasswordBody struct {
	ResetToken  string `json:"token" validate:"required,jwt"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
} // @Name ResetPassword

// @Tags Account
// @Summary Reset Password
// @Description Reset password
// @Accept json
// @Produce json
// @Param body body ResetPasswordBody true "Body"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/account/reset-password [post]
func ResetPassword(c *fiber.Ctx) error {
	body := ResetPasswordBody{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	claims, err := utils.GetJwtClaims(body.ResetToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: constants.InvalidJwt})
	}

	userResponse, err := crud.FindOneUserByEmail(claims.Subject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.UserNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	if !userResponse.IsActive {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: constants.UserInactive})
	}

	userUpdate := models.UserUpdate{
		Password: &body.NewPassword,
	}

	userResponse, err = crud.UpdateUser(userResponse.ID, userUpdate)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: constants.PasswordUpdated})
}
