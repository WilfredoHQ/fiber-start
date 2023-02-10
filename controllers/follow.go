package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
)

// @Tags Follows
// @Summary Read Follows
// @Description Read follows
// @Accept json
// @Produce json
// @Param followerId query string false "Follower Id"
// @Param followedId query string false "Followed Id"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} models.Follow
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/follows [get]
// @Security ApiKeyAuth
func ReadFollows(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	followerId := c.Query("followerId")
	followedId := c.Query("followedId")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)

	dbFollows, err := crud.FindAllFollows(followerId, followedId, int64(skip), int64(limit))
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(dbFollows)
}

// @Tags Follows
// @Summary Create Follow
// @Description Create follow
// @Accept json
// @Produce json
// @Param body body models.FollowCreate true "Body"
// @Success 201 {object} models.Follow
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/follows [post]
// @Security ApiKeyAuth
func CreateFollow(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	follow := models.Follow{
		FollowerID: currentUser.ID.Hex(),
	}

	if err := c.BodyParser(&follow); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Error{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&follow); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	dbFollows, err := crud.FindAllFollows(currentUser.ID.Hex(), follow.FollowedID, 0, 20)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	if len(dbFollows) > 0 {
		return c.Status(http.StatusCreated).JSON(dbFollows[0])
	}

	dbFollow, err := crud.InsertFollow(follow)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(dbFollow)
}

// @Tags Follows
// @Summary Read Follow
// @Description Read follow
// @Accept json
// @Produce json
// @Param follow_id path string true "Follow Id"
// @Success 200 {object} models.Follow
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/follows/{follow_id} [get]
// @Security ApiKeyAuth
func ReadFollow(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	followId := c.Params("followId")

	dbFollow, err := crud.FindOneFollowById(followId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "Follow not found"})
	}

	return c.Status(http.StatusOK).JSON(dbFollow)
}

// @Tags Follows
// @Summary Delete Follow
// @Description Delete follow
// @Accept json
// @Produce json
// @Param follow_id path string true "Follow Id"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/follows/{follow_id} [delete]
// @Security ApiKeyAuth
func DeleteFollow(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	followId := c.Params("followId")

	dbFollow, err := crud.FindOneFollowById(followId)
	if err != nil {
		return c.Status(http.StatusOK).JSON(models.Msg{Msg: "Follow deleted successfully"})
	}

	if !currentUser.IsSuperuser || dbFollow.FollowerID != currentUser.ID.Hex() {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: "The user does not have sufficient privileges"})
	}

	if err := crud.DeleteFollow(followId); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: "Follow deleted successfully"})
}
