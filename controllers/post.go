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

// @Tags Posts
// @Summary Get Posts
// @Description Get posts
// @Accept json
// @Produce json
// @Param userId query string false "User id"
// @Param search query string false "Search"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20) minimum(1)
// @Success 200 {array} models.PostResponse
// @Failure default {object} models.Error
// @Router /api/v1/posts [get]
// @Security ApiKeyAuth
func GetPosts(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	query := struct {
		UserID primitive.ObjectID `query:"userId"`
		Search string             `query:"search"`
		Skip   int                `query:"skip"`
		Limit  int                `query:"limit" validate:"min=1"`
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

	postsResponse, err := crud.FindAllPosts(query.UserID, query.Search, int64(query.Skip), int64(query.Limit))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(postsResponse)
}

// @Tags Posts
// @Summary Create Post
// @Description Create post
// @Accept json
// @Produce json
// @Param body body models.PostCreate true "Body"
// @Success 201 {object} models.PostResponse
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/posts [post]
// @Security ApiKeyAuth
func CreatePost(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	body := models.PostCreate{
		UserID: currentUser.ID,
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	postResponse, err := crud.InsertPost(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusCreated).JSON(postResponse)
}

// @Tags Posts
// @Summary Get Home Posts
// @Description Get home posts
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20) minimum(1)
// @Success 200 {array} models.PostResponse
// @Failure default {object} models.Error
// @Router /api/v1/posts/home [get]
// @Security ApiKeyAuth
func GetHomePosts(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	query := struct {
		Search string `query:"search"`
		Skip   int    `query:"skip"`
		Limit  int    `query:"limit" validate:"min=1"`
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

	postsResponse, err := crud.FindHomePosts(currentUser.ID, query.Search, int64(query.Skip), int64(query.Limit))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(postsResponse)
}

// @Tags Posts
// @Summary Get Post
// @Description Get post
// @Accept json
// @Produce json
// @Param post_id path string true "Post id"
// @Success 200 {object} models.PostResponse
// @Failure default {object} models.Error
// @Router /api/v1/posts/{post_id} [get]
// @Security ApiKeyAuth
func GetPost(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	params := struct {
		PostID primitive.ObjectID `params:"postId"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	postResponse, err := crud.FindOnePostById(params.PostID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.PostNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	return c.Status(http.StatusOK).JSON(postResponse)
}

// @Tags Posts
// @Summary Delete Post
// @Description Delete post
// @Accept json
// @Produce json
// @Param post_id path string true "Post id"
// @Success 200 {object} models.Msg
// @Failure default {object} models.Error
// @Router /api/v1/posts/{post_id} [delete]
// @Security ApiKeyAuth
func DeletePost(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	params := struct {
		PostID primitive.ObjectID `params:"postId"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	postResponse, err := crud.FindOnePostById(params.PostID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.PostNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	if postResponse.UserID != currentUser.ID && !currentUser.IsSuperuser {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: constants.InsufficientPrivileges})
	}

	if err := crud.DeletePost(params.PostID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: constants.PostDeleted})
}

// @Tags Posts
// @Summary Update Post
// @Description Update post
// @Accept json
// @Produce json
// @Param post_id path string true "Post id"
// @Param body body models.PostUpdate true "Body"
// @Success 200 {object} models.PostResponse
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error
// @Router /api/v1/posts/{post_id} [patch]
// @Security ApiKeyAuth
func UpdatePost(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	params := struct {
		PostID primitive.ObjectID `params:"postId"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	postResponse, err := crud.FindOnePostById(params.PostID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.PostNotFound})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
		}
	}

	if postResponse.UserID != currentUser.ID && !currentUser.IsSuperuser {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: constants.InsufficientPrivileges})
	}

	body := models.PostUpdate{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	postResponse, err = crud.UpdatePost(params.PostID, body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Error{Detail: constants.InternalServerError})
	}

	return c.Status(http.StatusOK).JSON(postResponse)
}
