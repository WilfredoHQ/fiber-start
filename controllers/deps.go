package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/crud"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
)

func currentUser(c *fiber.Ctx) (models.User, *fiber.Error) {
	claims := utils.GetLocalJwtClaims(c)

	dbUser, err := crud.FindOneUserById(claims.Subject)
	if err != nil {
		return dbUser, fiber.NewError(http.StatusNotFound, "User not found")
	}

	return dbUser, nil
}

func currentActiveUser(c *fiber.Ctx) (models.User, *fiber.Error) {
	dbUser, err := currentUser(c)
	if err != nil {
		return dbUser, err
	}

	if !dbUser.IsActive {
		return dbUser, fiber.NewError(http.StatusForbidden, "The user is inactive")
	}

	return dbUser, nil
}

func currentActiveSuperuser(c *fiber.Ctx) (models.User, *fiber.Error) {
	dbUser, err := currentActiveUser(c)
	if err != nil {
		return dbUser, err
	}

	if !dbUser.IsSuperuser {
		return dbUser, fiber.NewError(http.StatusForbidden, "The user does not have sufficient privileges")
	}

	return dbUser, nil
}
