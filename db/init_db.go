package db

import (
	"context"
	"time"

	"github.com/wilfredohq/fiber-start/config"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createSuperuser(ctx context.Context, client *mongo.Client) error {
	userCollection := GetCollection(client, "users")

	hashedPassword, err := utils.GetPasswordHash(config.Config.FirstSuperuserPassword)
	if err != nil {
		return err
	}

	fullName := "Superuser"
	isActive := true
	isSuperuser := true

	userCreate := models.UserCreate{
		FullName:    &fullName,
		Email:       &config.Config.FirstSuperuser,
		Password:    &hashedPassword,
		IsActive:    &isActive,
		IsSuperuser: &isSuperuser,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	filter := bson.M{"email": userCreate.Email}

	if result := userCollection.FindOne(ctx, filter); result.Err() == nil {
		return nil
	}

	if _, err := userCollection.InsertOne(ctx, userCreate); err != nil {
		return err
	}

	return nil
}
