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
type UpdateUserRequest struct {
	ID        uint   `json:"-"`
	Email     string `json:"email" validate:"omitempty,email"`
	FirstName string `json:"firstName" validate:"omitempty,gt=2"`
	LastName  string `json:"lastName" validate:"omitempty,gt=2"`
	Password  string `json:"password" validate:"omitempty,gt=6"`
}

type UsersResponse []*UserResponse

func NewUserResponse(user entity.User) *UserResponse {
	return &UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
func NewUsersResponse(users []entity.User) *UsersResponse {
	var usersResponse UsersResponse

	for _, user := range users {
		userResponse := UserResponse{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
		usersResponse = append(usersResponse, &userResponse)
	}
	return &usersResponse
}
func (c *CreateUserRequest) ToEntity() *entity.User {
	return &entity.User{
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Password:  c.Password,
	}
}
func (c *UpdateUserRequest) ToEntity() *entity.User {
	return &entity.User{
		ID:        c.ID,
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Password:  c.Password,
	}
}
