package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type PostUser struct {
	FullName  string `bson:"fullName" json:"fullName" validate:"required"`
	AvatarUrl string `bson:"avatarUrl" json:"avatarUrl" validate:"required"`
} // @Name PostUser

type PostResponse struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`
	Content   string             `bson:"content" json:"content" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt" validate:"required"`
	User      PostUser           `bson:"user" json:"user" validate:"required"`
} // @Name Post

type PostCreate struct {
	UserID    primitive.ObjectID `bson:"userId" swaggerignore:"true"`
	Content   *string            `bson:"content,omitempty" json:"content" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" swaggerignore:"true"`
	UpdatedAt time.Time          `bson:"updatedAt" swaggerignore:"true"`
} // @Name PostCreate

type PostUpdate struct {
	Content   *string   `bson:"content,omitempty" json:"content"`
	UpdatedAt time.Time `bson:"updatedAt" swaggerignore:"true"`
} // @Name PostUpdate
