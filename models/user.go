package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	FullName       string             `bson:"fullName"`
	Biography      string             `bson:"biography"`
	Location       string             `bson:"location"`
	Birthdate      time.Time          `bson:"birthdate"`
	Gender         string             `bson:"gender"`
	AvatarUrl      string             `bson:"avatarUrl"`
	CoverUrl       string             `bson:"coverUrl"`
	Email          string             `bson:"email"`
	Password       string             `bson:"password,omitempty"`
	IsActive       bool               `bson:"isActive"`
	IsSuperuser    bool               `bson:"isSuperuser"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	FollowersCount int                `bson:"followersCount"`
	FollowingCount int                `bson:"followingCount"`
}

type UserResponse struct {
	ID             primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	FullName       string             `bson:"fullName" json:"fullName" validate:"required"`
	Biography      string             `bson:"biography" json:"biography" validate:"required"`
	Location       string             `bson:"location" json:"location" validate:"required"`
	Birthdate      time.Time          `bson:"birthdate" json:"birthdate" validate:"required"`
	Gender         string             `bson:"gender" json:"gender" validate:"required"`
	AvatarUrl      string             `bson:"avatarUrl" json:"avatarUrl" validate:"required"`
	CoverUrl       string             `bson:"coverUrl" json:"coverUrl" validate:"required"`
	Email          string             `bson:"email" json:"email" validate:"required"`
	IsActive       bool               `bson:"isActive" json:"isActive" validate:"required"`
	IsSuperuser    bool               `bson:"isSuperuser" json:"isSuperuser" validate:"required"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt" validate:"required"`
	FollowersCount int                `bson:"followersCount" json:"followersCount" validate:"required"`
	FollowingCount int                `bson:"followingCount" json:"followingCount" validate:"required"`
} // @Name User

type UserCreate struct {
	FullName    *string    `bson:"fullName,omitempty" json:"fullName" validate:"required,min=3"`
	Biography   *string    `bson:"biography,omitempty" json:"biography"`
	Location    *string    `bson:"location,omitempty" json:"location"`
	Birthdate   *time.Time `bson:"birthdate,omitempty" json:"birthdate"`
	Gender      *string    `bson:"gender,omitempty" json:"gender"`
	AvatarUrl   *string    `bson:"avatarUrl,omitempty" json:"avatarUrl" validate:"omitempty,url"`
	CoverUrl    *string    `bson:"coverUrl,omitempty" json:"coverUrl" validate:"omitempty,url"`
	Email       *string    `bson:"email,omitempty" json:"email" validate:"required,email"`
	Password    *string    `bson:"password,omitempty" json:"password" validate:"required,min=8"`
	IsActive    *bool      `bson:"isActive,omitempty" json:"isActive"`
	IsSuperuser *bool      `bson:"isSuperuser,omitempty" json:"isSuperuser"`
	CreatedAt   time.Time  `bson:"createdAt" swaggerignore:"true"`
	UpdatedAt   time.Time  `bson:"updatedAt" swaggerignore:"true"`
} // @Name UserCreate

type UserUpdate struct {
	FullName    *string    `bson:"fullName,omitempty" json:"fullName" validate:"omitempty,min=3"`
	Biography   *string    `bson:"biography,omitempty" json:"biography"`
	Location    *string    `bson:"location,omitempty" json:"location"`
	Birthdate   *time.Time `bson:"birthdate,omitempty" json:"birthdate"`
	Gender      *string    `bson:"gender,omitempty" json:"gender"`
	AvatarUrl   *string    `bson:"avatarUrl,omitempty" json:"avatarUrl" validate:"omitempty,url"`
	CoverUrl    *string    `bson:"coverUrl,omitempty" json:"coverUrl" validate:"omitempty,url"`
	Password    *string    `bson:"password,omitempty" json:"password" validate:"omitempty,min=8"`
	IsActive    *bool      `bson:"isActive,omitempty" json:"isActive"`
	IsSuperuser *bool      `bson:"isSuperuser,omitempty" json:"isSuperuser"`
	UpdatedAt   time.Time  `bson:"updatedAt" swaggerignore:"true"`
} // @Name UserUpdate
