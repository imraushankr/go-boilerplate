package bank

import (
	"errors"
	"fmt"
	"time"
)

type Account struct {
	AccountNumber string
	HolderName    string
	Balance       float64
	AccountType   string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AccountService interface {
	Create(account Account) (*Account, error)
	Deposit(accountNumber string, amount float64) error
	Withdraw(accountNumber string, amount float64) error
	GetBalance(accountNumber string) (float64, error)
	GetAccountDetails(accountNumber string) (*Account, error)
	CloseAccount(accountNumber string) error
}

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrAccountExists     = errors.New("account already exists")
)

type accountService struct {
	accounts map[string]Account
}

func NewAccountService() AccountService {
	return &accountService{
		accounts: make(map[string]Account),
	}
}

func (ac *accountService) Create(account Account) (*Account, error) {
	if account.AccountNumber == "" || account.HolderName == "" || account.AccountType == "" {
		return nil, ErrInvalidInput
	}

	if _, exists := ac.accounts[account.AccountNumber]; exists {
		return nil, ErrAccountExists
	}

	account.Balance = 0.0
	account.Status = "Active"
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	ac.accounts[account.AccountNumber] = account
	return &account, nil
}

func (ac *accountService) Deposit(accountNumber string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	account, exists := ac.accounts[accountNumber]
	if !exists {
		return ErrAccountNotFound
	}

	if account.Status != "Active" {
		return fmt.Errorf("cannot deposit to an account with status: %s", account.Status)
	}

	account.Balance += amount
	account.UpdatedAt = time.Now()
	ac.accounts[accountNumber] = account

	return nil
}

func (ac *accountService) Withdraw(accountNumber string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	account, exists := ac.accounts[accountNumber]
	if !exists {
		return ErrAccountNotFound
	}

	if account.Status != "Active" {
		return fmt.Errorf("cannot withdraw from an account with status: %s", account.Status)
	}

	if account.Balance < amount {
		return ErrInsufficientFunds
	}

	account.Balance -= amount
	account.UpdatedAt = time.Now()
	ac.accounts[accountNumber] = account

	return nil
}

func (ac *accountService) GetBalance(accountNumber string) (float64, error) {
	account, exists := ac.accounts[accountNumber]
	if !exists {
		return 0, ErrAccountNotFound
	}

	return account.Balance, nil
}

func (ac *accountService) GetAccountDetails(accountNumber string) (*Account, error) {
	account, exists := ac.accounts[accountNumber]
	if !exists {
		return nil, ErrAccountNotFound
	}

	return &account, nil
}

func (ac *accountService) CloseAccount(accountNumber string) error {
	account, exists := ac.accounts[accountNumber]
	if !exists {
		return ErrAccountNotFound
	}

	if account.Balance != 0 {
		return errors.New("account balance must be zero to close the account")
	}

	account.Status = "Closed"
	account.UpdatedAt = time.Now()
	ac.accounts[accountNumber] = account

	return nil
}

func (a Account) DisplayAccountInfo() {
	fmt.Println("=== Account Information ===")
	fmt.Printf("Account Number: %s\n", a.AccountNumber)
	fmt.Printf("Holder Name: %s\n", a.HolderName)
	fmt.Printf("Balance: $%.2f\n", a.Balance)
	fmt.Printf("Type: %s\n", a.AccountType)
	fmt.Printf("Status: %s\n", a.Status)
	fmt.Printf("Created: %s\n", a.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", a.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("---------------------------")
}