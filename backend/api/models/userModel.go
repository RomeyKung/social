package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=5"`
	ImageUrl  string             `bson:"imageUrl" json:"imageUrl"`
	Bio       string             `bson:"bio" json:"bio"`
	Followers []string           `bson:"followers" json:"followers"`
	Following []string           `bson:"following" json:"following"`
}

type CreateUser struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=5"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	ImageUrl  string             `json:"imageUrl"`
	Bio       string             `json:"bio"`
	Followers []string           `json:"followers"`
	Following []string           `json:"following"`
}

func ToUserResponse(user UserModel) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		ImageUrl:  user.ImageUrl,
		Bio:       user.Bio,
		Followers: user.Followers,
		Following: user.Following,
	}
}
