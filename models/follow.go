package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FollowerID string             `bson:"followerId" json:"followerId"`
	FollowedID string             `bson:"followedId" json:"followedId" validate:"required"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
} // @Name Follow

type FollowCreate struct {
	FollowedID string `validate:"required"`
} // @Name FollowCreate

type FollowUpdate struct {
	FollowedID string
} // @Name FollowUpdate
