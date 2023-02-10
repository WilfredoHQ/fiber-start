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

func InsertPost(post models.Post) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	dbPost := models.Post{}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	result, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		return dbPost, err
	}

	dbPost = post
	dbPost.ID = result.InsertedID.(primitive.ObjectID)

	return dbPost, nil
}

func findOnePost(filter interface{}, opts ...*options.FindOneOptions) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	dbPost := models.Post{}

	if err := postCollection.FindOne(ctx, filter, opts...).Decode(&dbPost); err != nil {
		return dbPost, err
	}

	return dbPost, nil
}

func FindOnePostById(postId string) (models.Post, error) {
	postIdObj, _ := primitive.ObjectIDFromHex(postId)
	filter := bson.M{"_id": postIdObj}

	dbPost, err := findOnePost(filter)
	if err != nil {
		return dbPost, err
	}

	return dbPost, nil
}

func findPosts(filter interface{}, opts ...*options.FindOptions) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	dbPosts := []models.Post{}

	cur, err := postCollection.Find(ctx, filter, opts...)
	if err != nil {
		return dbPosts, err
	}

	if err := cur.All(ctx, &dbPosts); err != nil {
		return dbPosts, err
	}

	return dbPosts, nil
}

func FindAllPosts(userId string, search string, skip int64, limit int64) ([]models.Post, error) {
	filter := bson.M{"content": bson.M{"$regex": search, "$options": "i"}}
	if userId != "" {
		filter["userId"] = userId
	}
	opts := options.Find()
	opts.SetSort(bson.M{"createdAt": -1})
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	dbPosts, err := findPosts(filter, opts)
	if err != nil {
		return dbPosts, err
	}

	return dbPosts, nil
}

func aggregatePosts(pipeline interface{}, opts ...*options.AggregateOptions) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	dbPosts := []models.Post{}

	cur, err := postCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return dbPosts, err
	}

	if err := cur.All(ctx, &dbPosts); err != nil {
		return dbPosts, err
	}

	return dbPosts, nil
}

func AggregateHomePosts(followerId string, search string, skip int64, limit int64) ([]models.Post, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "follows",
			"localField":   "userId",
			"foreignField": "followedId",
			"as":           "followers",
		}},
		{"$match": bson.M{
			"content":              bson.M{"$regex": search, "$options": "i"},
			"followers.followerId": followerId,
		}},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	dbUsers, err := aggregatePosts(pipeline)
	if err != nil {
		return dbUsers, err
	}

	return dbUsers, nil
}

func UpdatePost(postId string, post models.Post) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	dbPost := models.Post{}

	post.UpdatedAt = time.Now()

	postIdObj, _ := primitive.ObjectIDFromHex(postId)
	filter := bson.M{"_id": postIdObj}
	update := bson.M{"$set": post}

	if _, err := postCollection.UpdateOne(ctx, filter, update); err != nil {
		return dbPost, err
	}

	dbPost = post

	return dbPost, nil
}

func DeletePost(postId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := db.GetCollection(db.DB, "posts")

	postIdObj, _ := primitive.ObjectIDFromHex(postId)
	filter := bson.M{"_id": postIdObj}

	if _, err := postCollection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
