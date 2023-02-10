package crud

import (
	"context"
	"time"

	"github.com/wilfredohq/fiber-start/db"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	dbUser := models.User{}

	user.Password = utils.GetPasswordHash(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return dbUser, err
	}

	dbUser = user
	dbUser.ID = result.InsertedID.(primitive.ObjectID)

	dbUser.Password = ""
	return dbUser, nil
}

func findOneUser(filter interface{}, opts ...*options.FindOneOptions) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	dbUser := models.User{}

	if err := userCollection.FindOne(ctx, filter, opts...).Decode(&dbUser); err != nil {
		return dbUser, err
	}

	dbUser.Password = ""
	return dbUser, nil
}

func FindOneUserById(userId string) (models.User, error) {
	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": userIdObj}

	dbUser, err := findOneUser(filter)
	if err != nil {
		return dbUser, err
	}

	return dbUser, nil
}

func FindOneUserByEmail(email string) (models.User, error) {
	filter := bson.M{"email": email}

	dbUser, err := findOneUser(filter)
	if err != nil {
		return dbUser, err
	}

	return dbUser, nil
}

func aggregateUsers(pipeline interface{}, opts ...*options.AggregateOptions) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	dbUsers := []models.User{}

	cur, err := userCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return dbUsers, err
	}

	for cur.Next(ctx) {
		dbUser := models.User{}

		if err := cur.Decode(&dbUser); err != nil {
			return dbUsers, err
		}

		dbUser.Password = ""
		dbUsers = append(dbUsers, dbUser)
	}

	return dbUsers, nil
}

func AggregateAllUsers(followerId string, followedId string, search string, skip int64, limit int64) ([]models.User, error) {
	match := bson.M{"fullName": bson.M{"$regex": search, "$options": "i"}}
	if followerId != "" {
		match["followers.followerId"] = followerId
	}
	if followedId != "" {
		match["following.followedId"] = followedId
	}
	pipeline := []bson.M{
		{"$addFields": bson.M{
			"userId": bson.M{"$toString": "$_id"},
		}},
		{"$lookup": bson.M{
			"from":         "follows",
			"localField":   "userId",
			"foreignField": "followerId",
			"as":           "following",
		}},
		{"$lookup": bson.M{
			"from":         "follows",
			"localField":   "userId",
			"foreignField": "followedId",
			"as":           "followers",
		}},
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	dbUsers, err := aggregateUsers(pipeline)
	if err != nil {
		return dbUsers, err
	}

	return dbUsers, nil
}

func UpdateUser(userId string, user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	dbUser := models.User{}

	if user.Password != "" {
		user.Password = utils.GetPasswordHash(user.Password)
	}
	user.UpdatedAt = time.Now()

	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": userIdObj}
	update := bson.M{"$set": user}

	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		return dbUser, err
	}

	dbUser = user

	dbUser.Password = ""
	return dbUser, nil
}

func DeleteUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": userIdObj}

	if _, err := userCollection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func AuthenticateUser(email string, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	dbUser := models.User{}

	filter := bson.M{"email": email}

	if err := userCollection.FindOne(ctx, filter).Decode(&dbUser); err != nil {
		return dbUser, err
	}

	if err := utils.VerifyPassword(password, dbUser.Password); err != nil {
		return dbUser, err
	}

	dbUser.Password = ""
	return dbUser, nil
}
