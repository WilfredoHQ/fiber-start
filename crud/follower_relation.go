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

func InsertFollowerRelation(followerRelationCreation models.FollowerRelationCreate) (models.FollowerRelationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := db.DB.StartSession()
	if err != nil {
		return models.FollowerRelationResponse{}, err
	}
	defer session.EndSession(ctx)

	followerRelationCollection := db.GetCollection(db.DB, "followerRelations")

	followerRelationCreation.CreatedAt = time.Now()
	followerRelationCreation.UpdatedAt = time.Now()

	transactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := followerRelationCollection.InsertOne(sessCtx, followerRelationCreation)
		if err != nil {
			return nil, err
		}

		if err := UpdateUserFollowersCount(sessCtx, *followerRelationCreation.FollowedID, 1); err != nil {
			return nil, err
		}

		if err := UpdateUserFollowingCount(sessCtx, followerRelationCreation.FollowerID, 1); err != nil {
			return nil, err
		}

		return result.InsertedID, nil
	}

	maxCommitTime := 10 * time.Second
	opts := options.Transaction().SetMaxCommitTime(&maxCommitTime)

	result, err := session.WithTransaction(ctx, transactionCallback, opts)
	if err != nil {
		return models.FollowerRelationResponse{}, err
	}

	return FindOneFollowerRelationById(result.(primitive.ObjectID))
}

func findOneFollowerRelation(filter interface{}, opts ...*options.FindOneOptions) (models.FollowerRelationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	followerRelationCollection := db.GetCollection(db.DB, "followerRelations")

	followerRelationResponse := models.FollowerRelationResponse{}

	if err := followerRelationCollection.FindOne(ctx, filter, opts...).Decode(&followerRelationResponse); err != nil {
		return models.FollowerRelationResponse{}, err
	}

	return followerRelationResponse, nil
}

func FindOneFollowerRelationById(followerRelationID primitive.ObjectID) (models.FollowerRelationResponse, error) {
	filter := bson.M{"_id": followerRelationID}

	return findOneFollowerRelation(filter)
}

func FindOneFollowerRelationByUserIds(followerID primitive.ObjectID, followedID primitive.ObjectID) (models.FollowerRelationResponse, error) {
	filter := bson.M{"followerId": followerID, "followedId": followedID}

	return findOneFollowerRelation(filter)
}

func DeleteFollowerRelation(followerRelationID primitive.ObjectID, followerRelationResponse models.FollowerRelationResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := db.DB.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	followerRelationCollection := db.GetCollection(db.DB, "followerRelations")

	filter := bson.M{"_id": followerRelationID}

	transactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if _, err := followerRelationCollection.DeleteOne(sessCtx, filter); err != nil {
			return nil, err
		}

		if err := UpdateUserFollowersCount(sessCtx, followerRelationResponse.FollowedID, -1); err != nil {
			return nil, err
		}

		if err := UpdateUserFollowingCount(sessCtx, followerRelationResponse.FollowerID, -1); err != nil {
			return nil, err
		}

		return nil, nil
	}

	maxCommitTime := 10 * time.Second
	opts := options.Transaction().SetMaxCommitTime(&maxCommitTime)

	if _, err := session.WithTransaction(ctx, transactionCallback, opts); err != nil {
		return err
	}

	return nil
}
