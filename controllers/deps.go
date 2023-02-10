package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/constants"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func currentUser(c *fiber.Ctx) (models.UserResponse, *fiber.Error) {
	userID := c.Locals("userId").(primitive.ObjectID)

	userResponse, err := crud.FindOneUserById(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return userResponse, fiber.NewError(http.StatusNotFound, constants.CurrentUserNotFound)
		} else {
			return userResponse, fiber.NewError(http.StatusInternalServerError, constants.InternalServerError)
		}
	}

	return userResponse, nil
}

func currentActiveUser(c *fiber.Ctx) (models.UserResponse, *fiber.Error) {
	userResponse, err := currentUser(c)
	if err != nil {
		return userResponse, err
	}

	if !userResponse.IsActive {
		return userResponse, fiber.NewError(http.StatusForbidden, constants.CurrentUserInactive)
	}

	return userResponse, nil
}

func currentActiveSuperuser(c *fiber.Ctx) (models.UserResponse, *fiber.Error) {
	userResponse, err := currentActiveUser(c)
	if err != nil {
		return userResponse, err
	}

	if !userResponse.IsSuperuser {
		return userResponse, fiber.NewError(http.StatusForbidden, constants.CurrentUserNotSuperuser)
	}

	return userResponse, nil
}
