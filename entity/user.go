package entity

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
}
type UserRepository interface {
	GetUser(ID uint) (*User, error)
	GetUsers() ([]User, error)
	CreateUser(User User) (*User, error)
	UpdateUser(User User) (*User, error)
	DeleteUser(ID uint) error
}
