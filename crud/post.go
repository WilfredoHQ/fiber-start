package crud

import (
	"context"
	"time"

	"github.com/wilfredohq/fiber-start/db"
	"github.com/wilfredohq/fiber-start/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertPost(postCreate models.PostCreate) (models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	postCreate.CreatedAt = time.Now()
	postCreate.UpdatedAt = time.Now()

	result, err := postCollection.InsertOne(ctx, postCreate)
	if err != nil {
		return models.PostResponse{}, err
	}

	return FindOnePostById(result.InsertedID.(primitive.ObjectID))
}

func findOnePost(match interface{}, opts ...*options.AggregateOptions) (models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": "$user"},
		{"$match": match},
	}

	cur, err := postCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return models.PostResponse{}, err
	}

	postResponse := models.PostResponse{}

	if cur.Next(ctx) {
		if err := cur.Decode(&postResponse); err != nil {
			return models.PostResponse{}, err
		}
	} else {
		return models.PostResponse{}, mongo.ErrNoDocuments
	}

	return postResponse, nil
}

func FindOnePostById(postID primitive.ObjectID) (models.PostResponse, error) {
	match := bson.M{"_id": postID}

	return findOnePost(match)
}

func findPosts(pipeline interface{}, opts ...*options.AggregateOptions) ([]models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	cur, err := postCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	postsResponse := []models.PostResponse{}

	if err := cur.All(ctx, &postsResponse); err != nil {
		return nil, err
	}

	return postsResponse, nil
}

func FindAllPosts(userID primitive.ObjectID, search string, skip int64, limit int64) ([]models.PostResponse, error) {
	match := bson.M{"content": bson.M{"$regex": search, "$options": "i"}}
	if !userID.IsZero() {
		match["userId"] = userID
	}

	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": "$user"},
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	return findPosts(pipeline)
}

func FindHomePosts(followerID primitive.ObjectID, search string, skip int64, limit int64) ([]models.PostResponse, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": "$user"},
		{"$lookup": bson.M{
			"from":         "followerRelations",
			"localField":   "userId",
			"foreignField": "followedId",
			"as":           "followers",
		}},
		{"$match": bson.M{
			"content":              bson.M{"$regex": search, "$options": "i"},
			"followers.followerId": followerID,
		}},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	return findPosts(pipeline)
}

func UpdatePost(postID primitive.ObjectID, postUpdate models.PostUpdate) (models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	postUpdate.UpdatedAt = time.Now()

	filter := bson.M{"_id": postID}
	update := bson.M{"$set": postUpdate}

	if _, err := postCollection.UpdateOne(ctx, filter, update); err != nil {
		return models.PostResponse{}, err
	}

	return FindOnePostById(postID)
}

func DeletePost(postID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	filter := bson.M{"_id": postID}

	if _, err := postCollection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
