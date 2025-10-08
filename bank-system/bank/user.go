package bank

import (
	"errors"
	"fmt"
)

type User struct {
	ID               int
	FirstName        string
	LastName         string
	Email            string
	Password         string
	Address          string
	Phone            string
	PanCardNumber    string
	AadharCardNumber string
	Accounts         []string
}

type UserService interface {
	Create(user User) (*User, error)
	Get(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user User) error
	Delete(id int) error
	List() ([]User, error)
	AddAccountToUser(userID int, accountNumber string) error
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
	if user.Email == "" || user.FirstName == "" || user.LastName == "" {
		return nil, ErrInvalidInput
	}

	for _, u := range us.users {
		if u.Email == user.Email {
			return nil, ErrEmailExists
		}
	}

	user.ID = us.nextID
	user.Accounts = []string{}
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

func (us *userService) GetByEmail(email string) (*User, error) {
	for _, user := range us.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, ErrUserNotFound
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

func (us *userService) AddAccountToUser(userID int, accountNumber string) error {
	user, exists := us.users[userID]
	if !exists {
		return ErrUserNotFound
	}

	for _, acc := range user.Accounts {
		if acc == accountNumber {
			return errors.New("account already linked to user")
		}
	}

	user.Accounts = append(user.Accounts, accountNumber)
	us.users[userID] = user
	return nil
}

func (u User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s %s, Email: %s, Phone: %s, Address: %s, Accounts: %v}",
		u.ID, u.FirstName, u.LastName, u.Email, u.Phone, u.Address, u.Accounts)
}

func (u User) DisplayUserInfo() {
	fmt.Println("=== User Information ===")
	fmt.Printf("ID: %d\n", u.ID)
	fmt.Printf("Name: %s %s\n", u.FirstName, u.LastName)
	fmt.Printf("Email: %s\n", u.Email)
	fmt.Printf("Phone: %s\n", u.Phone)
	fmt.Printf("Address: %s\n", u.Address)
	fmt.Printf("PAN Card: %s\n", u.PanCardNumber)
	fmt.Printf("Aadhar Card: %s\n", u.AadharCardNumber)
	fmt.Printf("Linked Accounts: %v\n", u.Accounts)
	fmt.Println("------------------------")
}