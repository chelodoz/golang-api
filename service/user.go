package service

import (
	"golang-api/dto"
	"golang-api/entity"
	"golang-api/util"
)

type UserService interface {
	CreateUser(user dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUsers() (*dto.UsersResponse, error)
	GetUser(ID uint) (*dto.UserResponse, error)
	UpdateUser(user dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ID uint) error
}

type userService struct {
	userRepository entity.UserRepository
}

func NewUserService(repository entity.UserRepository) UserService {
	return &userService{
		userRepository: repository,
	}
}

func (service *userService) CreateUser(createUserRequest dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashedPassword, err := util.HashPassword(createUserRequest.Password)
	if err != nil {
		return nil, err
	}
	createUserRequest.Password = hashedPassword

	user, err := service.userRepository.CreateUser(*createUserRequest.ToEntity())
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(*user), nil
}

func (service *userService) GetUser(ID uint) (*dto.UserResponse, error) {
	user, err := service.userRepository.GetUser(ID)
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(*user), nil
}
func (service *userService) GetUsers() (*dto.UsersResponse, error) {
	users, err := service.userRepository.GetUsers()
	if err != nil {
		return nil, err
	}
	return dto.NewUsersResponse(users), nil
}

func (service *userService) UpdateUser(updateUserRequest dto.UpdateUserRequest) (*dto.UserResponse, error) {

	user, err := service.userRepository.GetUser(updateUserRequest.ID)
	if err != nil {
		return nil, err
	}
	if updateUserRequest.Email != "" {
		user.Email = updateUserRequest.Email
	}
	if updateUserRequest.FirstName != "" {
		user.FirstName = updateUserRequest.FirstName
	}
	if updateUserRequest.LastName != "" {
		user.LastName = updateUserRequest.LastName
	}
	if updateUserRequest.Password != "" {
		user.Password = updateUserRequest.Password
	}

	user, err = service.userRepository.UpdateUser(*user)
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(*user), nil
}
func (service *userService) DeleteUser(ID uint) error {
	return service.userRepository.DeleteUser(ID)
}
