package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/constants"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Tags Follower relations
// @Summary Create Follower Relation
// @Description Create follower relation
// @Accept json
// @Produce json
// @Param body body models.FollowerRelationCreate true "Body"
// @Success 201 {object} models.FollowerRelationResponse
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/follower-relations [post]
// @Security ApiKeyAuth
func CreateFollowerRelation(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	body := models.FollowerRelationCreate{
		FollowerID: currentUser.ID,
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	if _, err := crud.FindOneFollowerRelationByUserIds(currentUser.ID, *body.FollowedID); err == nil {
		return c.Status(http.StatusConflict).JSON(models.Error{Detail: constants.FollowerRelationAlreadyRegistered})
	}

	followerRelationResponse, err := crud.InsertFollowerRelation(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	followerRelationResponse.HasData = true

	return c.Status(http.StatusCreated).JSON(followerRelationResponse)
}

// @Tags Follower relations
// @Summary Check Follower Relation
// @Description Check follower relation
// @Accept json
// @Produce json
// @Param user_id path string true "User id"
// @Success 200 {object} models.FollowerRelationResponse
// @Failure default {object} models.Error
// @Router /api/v1/follower-relations/following/{user_id} [get]
// @Security ApiKeyAuth
func CheckFollowerRelation(c *fiber.Ctx) error {
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

	followerRelationResponse, err := crud.FindOneFollowerRelationByUserIds(currentUser.ID, params.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusOK).JSON(models.FollowerRelationResponse{})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	followerRelationResponse.HasData = true

	return c.Status(http.StatusOK).JSON(followerRelationResponse)
}

// @Tags Follower relations
// @Summary Delete Follower Relation
// @Description Delete follower relation
// @Accept json
// @Produce json
// @Param follower_relation_id path string true "Follower relation id"
// @Success 200 {object} models.Msg
// @Failure default {object} models.Error
// @Router /api/v1/follower-relations/{follower_relation_id} [delete]
// @Security ApiKeyAuth
func DeleteFollowerRelation(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	params := struct {
		FollowerRelationID primitive.ObjectID `params:"followerRelationId"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	followerRelationResponse, err := crud.FindOneFollowerRelationById(params.FollowerRelationID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.FollowerRelationNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	if followerRelationResponse.FollowerID != currentUser.ID && !currentUser.IsSuperuser {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: constants.InsufficientPrivileges})
	}

	if err := crud.DeleteFollowerRelation(params.FollowerRelationID, followerRelationResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: constants.FollowerRelationDeleted})
}
