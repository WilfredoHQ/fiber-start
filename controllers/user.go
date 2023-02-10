package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/configs"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
)

// @Tags Users
// @Summary Read Users
// @Description Read users
// @Accept json
// @Produce json
// @Param followerId query string false "Follower Id"
// @Param followedId query string false "Followed Id"
// @Param search query string false "Search"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.User
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/users [get]
// @Security ApiKeyAuth
func ReadUsers(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	followerId := c.Query("followerId")
	followedId := c.Query("followedId")
	search := c.Query("search")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)

	dbUsers, err := crud.AggregateAllUsers(followerId, followedId, search, int64(skip), int64(limit))
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(dbUsers)
}

// @Tags Users
// @Summary Create User
// @Description Create user
// @Accept json
// @Produce json
// @Param body body models.UserCreate true "Body"
// @Success 201 {object} models.User
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/users [post]
func CreateUser(c *fiber.Ctx) error {
	if !configs.Env.UsersOpenRegistration {
		if _, fiberErr := currentActiveSuperuser(c); fiberErr != nil {
			return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
		}
	}

	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Error{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	if _, err := crud.FindOneUserByEmail(user.Email); err == nil {
		return c.Status(http.StatusConflict).JSON(models.Error{Detail: "The user with this email already exists in the system"})
	}

	if configs.Env.UsersOpenRegistration {
		user.IsActive = true
		user.IsSuperuser = false
	}

	dbUser, err := crud.InsertUser(user)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	utils.SendWelcomeEmail(dbUser.Email, dbUser.FullName)

	return c.Status(http.StatusCreated).JSON(dbUser)
}

// @Tags Users
// @Summary Read User
// @Description Read user
// @Accept json
// @Produce json
// @Param user_id path string true "User Id"
// @Success 200 {object} models.User
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/users/{user_id} [get]
// @Security ApiKeyAuth
func ReadUser(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	userId := c.Params("userId")

	dbUser, err := crud.FindOneUserById(userId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "User not found"})
	}

	return c.Status(http.StatusOK).JSON(dbUser)
}

// @Tags Users
// @Summary Delete User
// @Description Delete user
// @Accept json
// @Produce json
// @Param user_id path string true "User Id"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/users/{user_id} [delete]
// @Security ApiKeyAuth
func DeleteUser(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	userId := c.Params("userId")

	dbUser, err := crud.FindOneUserById(userId)
	if err != nil {
		return c.Status(http.StatusOK).JSON(models.Msg{Msg: "User deleted successfully"})
	}

	if !currentUser.IsSuperuser || dbUser.ID != currentUser.ID {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: "The user does not have sufficient privileges"})
	}

	if err := crud.DeleteUser(userId); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: "User deleted successfully"})
}

// @Tags Users
// @Summary Update User
// @Description Update user
// @Accept json
// @Produce json
// @Param user_id path string true "User Id"
// @Param body body models.UserUpdate true "Body"
// @Success 200 {object} models.User
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/users/{user_id} [patch]
// @Security ApiKeyAuth
func UpdateUser(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	userId := c.Params("userId")

	dbUser, err := crud.FindOneUserById(userId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "User not found"})
	}

	if !currentUser.IsSuperuser || dbUser.ID != currentUser.ID {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: "The user does not have sufficient privileges"})
	}

	user := dbUser

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Error{Detail: err.Error()})
	}

	exceptField := ""
	if user.Password == "" {
		exceptField = "Password"
	}

	validate := utils.NewValidator()
	if err := validate.StructExcept(&user, exceptField); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	user.ID = dbUser.ID
	user.Email = dbUser.Email
	if !currentUser.IsSuperuser {
		user.IsSuperuser = dbUser.IsSuperuser
	}
	user.CreatedAt = dbUser.CreatedAt

	dbUser, err = crud.UpdateUser(userId, user)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(dbUser)
}
