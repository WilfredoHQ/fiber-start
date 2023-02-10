package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId" json:"userId"`
	Content   string             `bson:"content" json:"content" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
} // @Name Post

type PostCreate struct {
	Content string `validate:"required"`
} // @Name PostCreate

type PostUpdate struct {
	Content string
} // @Name PostUpdate
