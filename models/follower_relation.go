package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowerRelation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FollowerID primitive.ObjectID `bson:"followerId"`
	FollowedID primitive.ObjectID `bson:"followedId"`
	CreatedAt  time.Time          `bson:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt"`
}

type FollowerRelationResponse struct {
	ID         primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	FollowerID primitive.ObjectID `bson:"followerId" json:"followerId" validate:"required"`
	FollowedID primitive.ObjectID `bson:"followedId" json:"followedId" validate:"required"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt" validate:"required"`
	HasData    bool               `json:"hasData" validate:"required"`
} // @Name FollowerRelation

type FollowerRelationCreate struct {
	FollowerID primitive.ObjectID  `bson:"followerId" swaggerignore:"true"`
	FollowedID *primitive.ObjectID `bson:"followedId,omitempty" json:"followedId" validate:"required"`
	CreatedAt  time.Time           `bson:"createdAt" swaggerignore:"true"`
	UpdatedAt  time.Time           `bson:"updatedAt" swaggerignore:"true"`
} // @Name FollowerRelationCreate
