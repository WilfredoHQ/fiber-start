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

func InsertUser(userCreate models.UserCreate) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	hashedPassword, err := utils.GetPasswordHash(*userCreate.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	userCreate.Password = &hashedPassword
	userCreate.CreatedAt = time.Now()
	userCreate.UpdatedAt = time.Now()

	result, err := userCollection.InsertOne(ctx, userCreate)
	if err != nil {
		return models.UserResponse{}, err
	}

	return FindOneUserById(result.InsertedID.(primitive.ObjectID))
}

func findOneUser(filter interface{}, opts ...*options.FindOneOptions) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	userResponse := models.UserResponse{}

	if err := userCollection.FindOne(ctx, filter, opts...).Decode(&userResponse); err != nil {
		return models.UserResponse{}, err
	}

	return userResponse, nil
}

func FindOneUserById(userID primitive.ObjectID) (models.UserResponse, error) {
	filter := bson.M{"_id": userID}

	return findOneUser(filter)
}

func FindOneUserByEmail(email string) (models.UserResponse, error) {
	filter := bson.M{"email": email}

	return findOneUser(filter)
}

func findUsers(pipeline interface{}, opts ...*options.AggregateOptions) ([]models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	cur, err := userCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	usersResponse := []models.UserResponse{}

	if err := cur.All(ctx, &usersResponse); err != nil {
		return nil, err
	}

	return usersResponse, nil
}

func FindAllUsers(followerID primitive.ObjectID, followedID primitive.ObjectID, search string, skip int64, limit int64) ([]models.UserResponse, error) {
	match := bson.M{"fullName": bson.M{"$regex": search, "$options": "i"}}
	if !followerID.IsZero() {
		match["followers.followerId"] = followerID
	}
	if !followedID.IsZero() {
		match["following.followedId"] = followedID
	}
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "followerRelations",
			"localField":   "_id",
			"foreignField": "followerId",
			"as":           "following",
		}},
		{"$lookup": bson.M{
			"from":         "followerRelations",
			"localField":   "_id",
			"foreignField": "followedId",
			"as":           "followers",
		}},
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	return findUsers(pipeline)
}

func UpdateUser(userID primitive.ObjectID, userUpdate models.UserUpdate) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	if userUpdate.Password != nil {
		hashedPassword, err := utils.GetPasswordHash(*userUpdate.Password)
		if err != nil {
			return models.UserResponse{}, err
		}

		userUpdate.Password = &hashedPassword
	}
	userUpdate.UpdatedAt = time.Now()

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": userUpdate}

	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		return models.UserResponse{}, err
	}

	return FindOneUserById(userID)
}

func updateUserCustomFields(ctx context.Context, userID primitive.ObjectID, update interface{}, opts ...*options.UpdateOptions) error {
	userCollection := db.GetCollection(db.DB, "users")

	filter := bson.M{"_id": userID}

	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func UpdateUserFollowersCount(ctx context.Context, userID primitive.ObjectID, followerCountDelta int) error {
	update := bson.M{"$inc": bson.M{"followersCount": followerCountDelta}}

	return updateUserCustomFields(ctx, userID, update)
}

func UpdateUserFollowingCount(ctx context.Context, userID primitive.ObjectID, followingCountDelta int) error {
	update := bson.M{"$inc": bson.M{"followingCount": followingCountDelta}}

	return updateUserCustomFields(ctx, userID, update)
}

func AuthenticateUser(email string, password string) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(db.DB, "users")

	dbUser := models.User{}

	filter := bson.M{"email": email}

	if err := userCollection.FindOne(ctx, filter).Decode(&dbUser); err != nil {
		return models.UserResponse{}, err
	}

	if err := utils.VerifyPassword(password, dbUser.Password); err != nil {
		return models.UserResponse{}, err
	}

	userResponse := models.UserResponse{
		ID:             dbUser.ID,
		FullName:       dbUser.FullName,
		Biography:      dbUser.Biography,
		Location:       dbUser.Location,
		Birthdate:      dbUser.Birthdate,
		Gender:         dbUser.Gender,
		AvatarUrl:      dbUser.AvatarUrl,
		CoverUrl:       dbUser.CoverUrl,
		Email:          dbUser.Email,
		IsActive:       dbUser.IsActive,
		IsSuperuser:    dbUser.IsSuperuser,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		FollowersCount: dbUser.FollowersCount,
		FollowingCount: dbUser.FollowingCount,
	}

	return userResponse, nil
}
