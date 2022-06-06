package service

import (
	"golang-api/entity"
	"golang-api/util"
)

type LoginService interface {
	Login(email string, password string) bool
}

type loginService struct {
	userRepository entity.UserRepository
}

func NewLoginService(repository entity.UserRepository) LoginService {
	return &loginService{
		userRepository: repository,
	}
}

func (service *loginService) Login(email string, password string) bool {
	user, err := service.userRepository.GetUserByEmail(email)
	if err != nil {
		return false
	}

	err = util.CheckPassword(password, user.Password)
	return err == nil
}
