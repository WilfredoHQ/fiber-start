package crud

import (
	"context"
	"time"

	"github.com/wilfredohq/fiber-start/db"
	"github.com/wilfredohq/fiber-start/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertFollow(follow models.Follow) (models.Follow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	followCollection := db.GetCollection(db.DB, "follows")

	dbFollow := models.Follow{}

	follow.CreatedAt = time.Now()
	follow.UpdatedAt = time.Now()

	result, err := followCollection.InsertOne(ctx, follow)
	if err != nil {
		return dbFollow, err
	}

	dbFollow = follow
	dbFollow.ID = result.InsertedID.(primitive.ObjectID)

	return dbFollow, nil
}

func findOneFollow(filter interface{}, opts ...*options.FindOneOptions) (models.Follow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	followCollection := db.GetCollection(db.DB, "follows")

	dbFollow := models.Follow{}

	if err := followCollection.FindOne(ctx, filter, opts...).Decode(&dbFollow); err != nil {
		return dbFollow, err
	}

	return dbFollow, nil
}

func FindOneFollowById(followId string) (models.Follow, error) {
	followIdObj, _ := primitive.ObjectIDFromHex(followId)
	filter := bson.M{"_id": followIdObj}

	dbFollow, err := findOneFollow(filter)
	if err != nil {
		return dbFollow, err
	}

	return dbFollow, nil
}

func findFollows(filter interface{}, opts ...*options.FindOptions) ([]models.Follow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	followCollection := db.GetCollection(db.DB, "follows")

	dbFollows := []models.Follow{}

	cur, err := followCollection.Find(ctx, filter, opts...)
	if err != nil {
		return dbFollows, err
	}

	if err := cur.All(ctx, &dbFollows); err != nil {
		return dbFollows, err
	}

	return dbFollows, nil
}

func FindAllFollows(followerId string, followedId string, skip int64, limit int64) ([]models.Follow, error) {
	filter := bson.M{}
	if followerId != "" {
		filter["followerId"] = followerId
	}
	if followedId != "" {
		filter["followedId"] = followedId
	}
	opts := options.Find()
	opts.SetSort(bson.M{"createdAt": -1})
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	dbFollows, err := findFollows(filter, opts)
	if err != nil {
		return dbFollows, err
	}

	return dbFollows, nil
}

func DeleteFollow(followId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	followCollection := db.GetCollection(db.DB, "follows")

	followIdObj, _ := primitive.ObjectIDFromHex(followId)
	filter := bson.M{"_id": followIdObj}

	if _, err := followCollection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
