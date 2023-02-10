package db

import (
	"context"
	"time"

	"github.com/wilfredohq/fiber-start/configs"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createSuperuser(ctx context.Context, client *mongo.Client) (models.User, error) {
	userCollection := GetCollection(client, "users")

	user := models.User{
		FullName:    "Superuser",
		Email:       configs.Env.FirstSuperuser,
		Password:    utils.GetPasswordHash(configs.Env.FirstSuperuserPassword),
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	dbUser := models.User{}

	filter := bson.M{"email": user.Email}

	if err := userCollection.FindOne(ctx, filter).Decode(&dbUser); err == nil {
		return dbUser, nil
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return dbUser, err
	}

	dbUser = user
	dbUser.ID = result.InsertedID.(primitive.ObjectID)

	dbUser.Password = ""
	return dbUser, nil
}
