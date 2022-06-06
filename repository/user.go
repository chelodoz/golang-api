package repository

import (
	"golang-api/entity"
	"time"

	"gorm.io/gorm"
)

type UserGorm struct {
	ID        uint      `gorm:"primary_key;auto_increment"`
	FirstName string    `gorm:"type:varchar(32)"`
	LastName  string    `gorm:"type:varchar(32)"`
	Email     string    `gorm:"type:varchar(256);UNIQUE"`
	Password  string    `gorm:"type:varchar(256)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (UserGorm) TableName() string {
	return "users"
}

func (u UserGorm) ToEntity() (*entity.User, error) {
	return &entity.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}, nil
}

func NewUserGorm(u entity.User) UserGorm {
	return UserGorm{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) entity.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (userRepository *userRepository) CreateUser(user entity.User) (*entity.User, error) {
	userGorm := NewUserGorm(user)
	err := userRepository.DB.Create(&userGorm).Error
	if err != nil {
		return nil, err
	}
	return userGorm.ToEntity()
}

func (userRepository *userRepository) UpdateUser(user entity.User) (*entity.User, error) {
	userGorm := NewUserGorm(user)
	err := userRepository.DB.Updates(&userGorm).Error
	if err != nil {
		return nil, err
	}
	return userGorm.ToEntity()
}

func (userRepository *userRepository) DeleteUser(ID uint) error {
	err := userRepository.DB.Delete(&UserGorm{}, ID).Error
	return err
}

func (userRepository *userRepository) GetUserByID(ID uint) (*entity.User, error) {
	userGorm := &UserGorm{ID: ID}
	err := userRepository.DB.First(&userGorm).Error
	if err != nil {
		return &entity.User{}, err
	}
	return userGorm.ToEntity()
}

func (userRepository *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	userGorm := &UserGorm{Email: email}
	err := userRepository.DB.First(&userGorm, "email = ?", email).Error
	if err != nil {
		return &entity.User{}, err
	}
	return userGorm.ToEntity()
}

func (userRepository *userRepository) GetUsers() ([]entity.User, error) {
	var usersGorm []UserGorm
	var users []entity.User
	err := userRepository.DB.Find(&usersGorm).Error

	for _, userGorm := range usersGorm {
		user, err := userGorm.ToEntity()
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, err
}
