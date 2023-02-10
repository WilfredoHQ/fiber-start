package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/config"
	"github.com/wilfredohq/fiber-start/constants"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Tags Users
// @Summary Get Users
// @Description Get users
// @Accept json
// @Produce json
// @Param followerId query string false "Follower id"
// @Param followedId query string false "Followed id"
// @Param search query string false "Search"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20) minimum(1)
// @Success 200 {array} models.UserResponse
// @Failure default {object} models.Error
// @Router /api/v1/users [get]
// @Security ApiKeyAuth
func GetUsers(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	query := struct {
		FollowerID primitive.ObjectID `query:"followerId"`
		FollowedID primitive.ObjectID `query:"followedId"`
		Search     string             `query:"search"`
		Skip       int                `query:"skip"`
		Limit      int                `query:"limit" validate:"min=1"`
	}{
		Limit: 20,
	}

	if err := c.QueryParser(&query); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&query); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	usersResponse, err := crud.FindAllUsers(query.FollowerID, query.FollowedID, query.Search, int64(query.Skip), int64(query.Limit))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(usersResponse)
}

// @Tags Users
// @Summary Create User
// @Description Create user
// @Accept json
// @Produce json
// @Param body body models.UserCreate true "Body"
// @Success 201 {object} models.UserResponse
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/users [post]
func CreateUser(c *fiber.Ctx) error {
	if !config.Config.UsersOpenRegistration {
		if _, fiberErr := currentActiveSuperuser(c); fiberErr != nil {
			return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
		}
	}

	body := models.UserCreate{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	if _, err := crud.FindOneUserByEmail(*body.Email); err == nil {
		return c.Status(http.StatusConflict).JSON(models.Error{Detail: constants.UserAlreadyRegistered})
	}

	if config.Config.UsersOpenRegistration {
		isActive := true
		isSuperuser := false

		body.IsActive = &isActive
		body.IsSuperuser = &isSuperuser
	}

	userResponse, err := crud.InsertUser(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	utils.SendWelcomeEmail(userResponse.Email, userResponse.FullName)

	return c.Status(http.StatusCreated).JSON(userResponse)
}

// @Tags Users
// @Summary Get User
// @Description Get user
// @Accept json
// @Produce json
// @Param user_id path string true "User id"
// @Success 200 {object} models.UserResponse
// @Failure default {object} models.Error
// @Router /api/v1/users/{user_id} [get]
// @Security ApiKeyAuth
func GetUser(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	params := struct {
		UserID primitive.ObjectID `params:"userId"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	userResponse, err := crud.FindOneUserById(params.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.UserNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	return c.Status(http.StatusOK).JSON(userResponse)
}

// @Tags Users
// @Summary Update User
// @Description Update user
// @Accept json
// @Produce json
// @Param user_id path string true "User id"
// @Param body body models.UserUpdate true "Body"
// @Success 200 {object} models.UserResponse
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/users/{user_id} [patch]
// @Security ApiKeyAuth
func UpdateUser(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	params := struct {
		UserID primitive.ObjectID `params:"userId"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	userResponse, err := crud.FindOneUserById(params.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.UserNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	if userResponse.ID != currentUser.ID && !currentUser.IsSuperuser {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: constants.InsufficientPrivileges})
	}

	body := models.UserUpdate{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	if !currentUser.IsSuperuser {
		body.IsSuperuser = &userResponse.IsSuperuser
	}

	userResponse, err = crud.UpdateUser(params.UserID, body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(userResponse)
}
