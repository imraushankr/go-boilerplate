package bank

import (
	"errors"
	"fmt"
)

type BankingSystem struct {
	users        UserService
	accounts     AccountService
	transactions TransactionService
}

func NewBankingSystem() *BankingSystem {
	bankingSystem := &BankingSystem{
		users:    NewUserService(),
		accounts: NewAccountService(),
	}

	bankingSystem.transactions = NewTransactionService(bankingSystem)
	return bankingSystem
}

func (bs *BankingSystem) CreateUser(firstName, lastName, email, password, address, phone, panCard, aadharCard string) (*User, error) {
	user := User{
		FirstName:        firstName,
		LastName:         lastName,
		Email:            email,
		Password:         password,
		Address:          address,
		Phone:            phone,
		PanCardNumber:    panCard,
		AadharCardNumber: aadharCard,
	}

	createdUser, err := bs.users.Create(user)
	if err != nil {
		return nil, err
	}

	fmt.Printf("User created successfully: %s %s (ID: %d)\n", firstName, lastName, createdUser.ID)
	return createdUser, nil
}

func (bs *BankingSystem) CreateAccount(accountNumber, holderName, accountType string, userID int) error {
	// Verify user exists
	user, err := bs.users.Get(userID)
	if err != nil {
		return errors.New("user not found")
	}

	account := Account{
		AccountNumber: accountNumber,
		HolderName:    holderName,
		AccountType:   accountType,
	}

	createdAccount, err := bs.accounts.Create(account)
	if err != nil {
		return err
	}

	// Link account to user
	err = bs.users.AddAccountToUser(userID, accountNumber)
	if err != nil {
		return err
	}

	fmt.Printf("Account created successfully for %s (User ID: %d)\n", holderName, userID)
	user.DisplayUserInfo()
	createdAccount.DisplayAccountInfo()
	return nil
}

func (bs *BankingSystem) Deposit(accountNumber string, amount float64) error {
	err := bs.accounts.Deposit(accountNumber, amount)
	if err != nil {
		return err
	}

	// Record transaction
	_, err = bs.transactions.CreateTransaction(Deposit, "", accountNumber, amount, "Cash deposit")
	if err != nil {
		fmt.Printf("Warning: Failed to record deposit transaction: %v\n", err)
	}

	fmt.Printf("Deposited $%.2f to account %s\n", amount, accountNumber)
	return nil
}

func (bs *BankingSystem) Withdraw(accountNumber string, amount float64) error {
	err := bs.accounts.Withdraw(accountNumber, amount)
	if err != nil {
		return err
	}

	// Record transaction
	_, err = bs.transactions.CreateTransaction(Withdrawal, accountNumber, "", amount, "Cash withdrawal")
	if err != nil {
		fmt.Printf("Warning: Failed to record withdrawal transaction: %v\n", err)
	}

	fmt.Printf("Withdrew $%.2f from account %s\n", amount, accountNumber)
	return nil
}

func (bs *BankingSystem) Transfer(fromAccount, toAccount string, amount float64) error {
	// Withdraw from source account
	err := bs.accounts.Withdraw(fromAccount, amount)
	if err != nil {
		return err
	}

	// Deposit to destination account
	err = bs.accounts.Deposit(toAccount, amount)
	if err != nil {
		// Rollback the withdrawal if deposit fails
		bs.accounts.Deposit(fromAccount, amount)
		return err
	}

	// Record transaction
	_, err = bs.transactions.CreateTransaction(Transfer, fromAccount, toAccount, amount, "Fund transfer")
	if err != nil {
		fmt.Printf("Warning: Failed to record transfer transaction: %v\n", err)
	}

	fmt.Printf("Transferred $%.2f from %s to %s\n", amount, fromAccount, toAccount)
	return nil
}

func (bs *BankingSystem) GetAccount(accountNumber string) (*Account, error) {
	return bs.accounts.GetAccountDetails(accountNumber)
}

func (bs *BankingSystem) GetUser(userID int) (*User, error) {
	return bs.users.Get(userID)
}

func (bs *BankingSystem) GetUserByEmail(email string) (*User, error) {
	return bs.users.GetByEmail(email)
}

func (bs *BankingSystem) GetBalance(accountNumber string) (float64, error) {
	return bs.accounts.GetBalance(accountNumber)
}

func (bs *BankingSystem) ListAllAccounts() {
	fmt.Println("\n=== All Accounts ===")
	users, _ := bs.users.List()
	accountFound := false
	for _, user := range users {
		for _, accountNumber := range user.Accounts {
			account, err := bs.accounts.GetAccountDetails(accountNumber)
			if err == nil {
				account.DisplayAccountInfo()
				accountFound = true
			}
		}
	}
	if !accountFound {
		fmt.Println("No accounts found.")
	}
}

func (bs *BankingSystem) ListAllUsers() {
	fmt.Println("\n=== All Users ===")
	users, err := bs.users.List()
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}

	for _, user := range users {
		user.DisplayUserInfo()
	}
}

func (bs *BankingSystem) GetTransactionHistory() {
	fmt.Println("\n=== Transaction History ===")
	transactions := bs.transactions.GetAllTransactions()
	DisplayTransactions(transactions, "All Transactions")
}

func (bs *BankingSystem) GetUserAccounts(userID int) {
	user, err := bs.users.Get(userID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\n=== Accounts for User %d: %s %s ===\n", user.ID, user.FirstName, user.LastName)
	if len(user.Accounts) == 0 {
		fmt.Println("No accounts found for this user.")
		return
	}
	
	for _, accountNumber := range user.Accounts {
		account, err := bs.accounts.GetAccountDetails(accountNumber)
		if err == nil {
			account.DisplayAccountInfo()
		}
	}
}

func (bs *BankingSystem) GetAccountTransactions(accountNumber string) {
	transactions, err := bs.transactions.GetTransactionsByAccount(accountNumber)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	DisplayTransactions(transactions, fmt.Sprintf("Transactions for Account %s", accountNumber))
}

func (bs *BankingSystem) GetTransactionSummary(accountNumber string) {
	summary := bs.transactions.GetTransactionSummary(accountNumber)
	summary.DisplayTransactionSummary()
}

func (bs *BankingSystem) CloseAccount(accountNumber string) error {
	return bs.accounts.CloseAccount(accountNumber)
}