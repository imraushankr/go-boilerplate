package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID              int
	FirstName       string
	LastName        string
	Email           string
	Password        string
	Address         string
	Phone           string
	PanCardNumber   string
	AddarCardNumber string
}

type UserService interface {
	Create(user User) (*User, error)
	Get(id int) (*User, error)
	Update(user User) error
	Delete(id int) error
	List() ([]User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
)

type userService struct {
	users  map[int]User
	nextID int
}

func NewUserService() UserService {
	return &userService{
		users:  make(map[int]User),
		nextID: 1,
	}
}

func (us *userService) Create(user User) (*User, error) {
	if user.Email == "" {
		return nil, ErrInvalidInput
	}

	for _, u := range us.users {
		if u.Email == user.Email {
			return nil, ErrEmailExists
		}
	}

	user.ID = us.nextID
	us.users[us.nextID] = user
	us.nextID++

	return &user, nil
}

func (us *userService) Get(id int) (*User, error) {
	user, exists := us.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (us *userService) Update(user User) error {
	_, exists := us.users[user.ID]
	if !exists {
		return ErrUserNotFound
	}

	us.users[user.ID] = user
	return nil
}

func (us *userService) Delete(id int) error {
	_, exists := us.users[id]
	if !exists {
		return ErrUserNotFound
	}

	delete(us.users, id)
	return nil
}

func (us *userService) List() ([]User, error) {
	userList := make([]User, 0, len(us.users))

	for _, user := range us.users {
		userList = append(userList, user)
	}

	return userList, nil
}

func (u User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s %s, Email: %s, Phone: %s, Address: %s}",
		u.ID, u.FirstName, u.LastName, u.Email, u.Phone, u.Address)
}

func main() {
	usrService := NewUserService()

	usr := User{
		FirstName:       "Raushan",
		LastName:        "Kumar",
		Email:           "raushan.kumar@hk.com",
		Password:        "RK@tr",
		Address:         "Motihari, Bihar",
		Phone:           "7645927364",
		PanCardNumber:   "DFFGD7657JKHG",
		AddarCardNumber: "3245733353443626",
	}

	createdUsr, err := usrService.Create(usr)
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return
	}

	fmt.Printf("New created User :: %s\n", createdUsr)
}