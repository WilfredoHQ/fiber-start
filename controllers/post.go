package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
)

// @Tags Posts
// @Summary Read Posts
// @Description Read posts
// @Accept json
// @Produce json
// @Param userId query string false "User Id"
// @Param search query string false "Search"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} models.Post
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/posts [get]
// @Security ApiKeyAuth
func ReadPosts(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	userId := c.Query("userId")
	search := c.Query("search")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)

	dbPosts, err := crud.FindAllPosts(userId, search, int64(skip), int64(limit))
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(dbPosts)
}

// @Tags Posts
// @Summary Create Post
// @Description Create post
// @Accept json
// @Produce json
// @Param body body models.PostCreate true "Body"
// @Success 201 {object} models.Post
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/posts [post]
// @Security ApiKeyAuth
func CreatePost(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	post := models.Post{
		UserID: currentUser.ID.Hex(),
	}

	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Error{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&post); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	dbPost, err := crud.InsertPost(post)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(dbPost)
}

// @Tags Posts
// @Summary Read Home Posts
// @Description Read home posts
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} models.Post
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/posts/home [get]
// @Security ApiKeyAuth
func ReadHomePosts(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	search := c.Query("search")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)

	dbPosts, err := crud.AggregateHomePosts(currentUser.ID.Hex(), search, int64(skip), int64(limit))
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(dbPosts)
}

// @Tags Posts
// @Summary Read Post
// @Description Read post
// @Accept json
// @Produce json
// @Param post_id path string true "Post Id"
// @Success 200 {object} models.Post
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/posts/{post_id} [get]
// @Security ApiKeyAuth
func ReadPost(c *fiber.Ctx) error {
	if _, fiberErr := currentActiveUser(c); fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	postId := c.Params("postId")

	dbPost, err := crud.FindOnePostById(postId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "Post not found"})
	}

	return c.Status(http.StatusOK).JSON(dbPost)
}

// @Tags Posts
// @Summary Delete Post
// @Description Delete post
// @Accept json
// @Produce json
// @Param post_id path string true "Post Id"
// @Success 200 {object} models.Msg
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/posts/{post_id} [delete]
// @Security ApiKeyAuth
func DeletePost(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	postId := c.Params("postId")

	dbPost, err := crud.FindOnePostById(postId)
	if err != nil {
		return c.Status(http.StatusOK).JSON(models.Msg{Msg: "Post deleted successfully"})
	}

	if !currentUser.IsSuperuser || dbPost.UserID != currentUser.ID.Hex() {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: "The user does not have sufficient privileges"})
	}

	if err := crud.DeletePost(postId); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.Msg{Msg: "Post deleted successfully"})
}

// @Tags Posts
// @Summary Update Post
// @Description Update post
// @Accept json
// @Produce json
// @Param post_id path string true "Post Id"
// @Param body body models.PostUpdate true "Body"
// @Success 200 {object} models.Post
// @Failure 422 {object} models.ValidationError
// @Failure default {object} models.Error "Error"
// @Router /api/v1/posts/{post_id} [patch]
// @Security ApiKeyAuth
func UpdatePost(c *fiber.Ctx) error {
	currentUser, fiberErr := currentActiveUser(c)
	if fiberErr != nil {
		return c.Status(fiberErr.Code).JSON(models.Error{Detail: fiberErr.Message})
	}

	postId := c.Params("postId")

	dbPost, err := crud.FindOnePostById(postId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: "Post not found"})
	}

	if !currentUser.IsSuperuser || dbPost.UserID != currentUser.ID.Hex() {
		return c.Status(http.StatusForbidden).JSON(models.Error{Detail: "The user does not have sufficient privileges"})
	}

	post := dbPost

	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Error{Detail: err.Error()})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&post); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(models.ValidationError{Detail: utils.ValidatorErrors(err)})
	}

	post.ID = dbPost.ID
	post.UserID = dbPost.UserID
	post.CreatedAt = dbPost.CreatedAt

	dbPost, err = crud.UpdatePost(postId, post)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(dbPost)
}
