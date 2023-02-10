package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName    string             `bson:"fullName" json:"fullName"`
	Email       string             `bson:"email" json:"email" validate:"required,email"`
	Password    string             `bson:"password,omitempty" json:"password,omitempty" validate:"required,gte=8"`
	IsActive    bool               `bson:"isActive" json:"isActive" validate:"boolean"`
	IsSuperuser bool               `bson:"isSuperuser" json:"isSuperuser" validate:"boolean"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
} // @Name User

type UserCreate struct {
	FullName    string
	Email       string `validate:"required"`
	Password    string `validate:"required,gte=8"`
	IsActive    bool
	IsSuperuser bool
} // @Name UserCreate

type UserUpdate struct {
	FullName    string
	Password    string `validate:"gte=8"`
	IsActive    bool
	IsSuperuser bool
} // @Name UserUpdate
