package dto

import "golang-api/entity"

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required,gt=6"`
}

type UserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func NewUserResponse(user entity.User) *UserResponse {
	return &UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
func (c *CreateUserRequest) ToEntity() *entity.User {
	return &entity.User{
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
}
