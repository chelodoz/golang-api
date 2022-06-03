package service

import (
	"golang-api/dto"
	"golang-api/entity"
)

type UserService interface {
	CreateUser(user dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUsers() (*dto.UsersResponse, error)
	// GetUser() (entity.User, error)
	// UpdateUser(User entity.User) (*entity.User, error)
	// DeleteUser(ID uint) error
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

	user, err := service.userRepository.CreateUser(*createUserRequest.ToEntity())

	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(*user), nil
}

// func (service *userService) GetUser(user entity.User) (entity.User, error) {
// 	return service.userRepository.GetUser(user)
// }
func (service *userService) GetUsers() (*dto.UsersResponse, error) {
	users, err := service.userRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	return dto.NewUsersResponse(users), nil
}

// func (service *userService) UpdateUser(user entity.User) (entity.User, error) {
// 	return service.userRepository.UpdateUser(user)
// }
// func (service *userService) DeleteUser(user entity.User) (entity.User, error) {
// 	return service.userRepository.DeleteUser(user)
// }
