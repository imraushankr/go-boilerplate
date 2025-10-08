package bank

import (
	"errors"
	"fmt"
	"time"
)

type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
	Interest   TransactionType = "INTEREST"
	Fee        TransactionType = "FEE"
)

type TransactionStatus string

const (
	Pending   TransactionStatus = "PENDING"
	Completed TransactionStatus = "COMPLETED"
	Failed    TransactionStatus = "FAILED"
	Cancelled TransactionStatus = "CANCELLED"
)

type Transaction struct {
	ID              string
	Type            TransactionType
	Status          TransactionStatus
	FromAccount     string
	ToAccount       string
	Amount          float64
	Timestamp       time.Time
	Description     string
	ReferenceNumber string
	BalanceAfter    float64
	Fee             float64
}

type TransactionService interface {
	CreateTransaction(tType TransactionType, fromAcc, toAcc string, amount float64, description string) (*Transaction, error)
	GetTransaction(transactionID string) (*Transaction, error)
	GetTransactionsByAccount(accountNumber string) ([]*Transaction, error)
	GetTransactionsByUser(userID int) ([]*Transaction, error)
	GetTransactionsByType(tType TransactionType) ([]*Transaction, error)
	GetTransactionsByDateRange(startDate, endDate time.Time) ([]*Transaction, error)
	UpdateTransactionStatus(transactionID string, status TransactionStatus) error
	GetAllTransactions() []*Transaction
	GetTransactionSummary(accountNumber string) *TransactionSummary
}

type TransactionSummary struct {
	AccountNumber     string
	TotalDeposits     float64
	TotalWithdrawals  float64
	TotalTransfersOut float64
	TotalTransfersIn  float64
	TotalFees         float64
	TransactionCount  int
	LastTransaction   time.Time
}

type transactionService struct {
	transactions  map[string]*Transaction
	bankingSystem *BankingSystem
}

func NewTransactionService(bankingSystem *BankingSystem) TransactionService {
	return &transactionService{
		transactions:  make(map[string]*Transaction),
		bankingSystem: bankingSystem,
	}
}

func generateTransactionID() string {
	return fmt.Sprintf("TXN%d", time.Now().UnixNano())
}

func generateReferenceNumber() string {
	return fmt.Sprintf("REF%d", time.Now().UnixNano())
}

func (ts *transactionService) CreateTransaction(tType TransactionType, fromAcc, toAcc string, amount float64, description string) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("transaction amount must be positive")
	}

	transaction := &Transaction{
		ID:              generateTransactionID(),
		Type:            tType,
		Status:          Pending,
		FromAccount:     fromAcc,
		ToAccount:       toAcc,
		Amount:          amount,
		Timestamp:       time.Now(),
		Description:     description,
		ReferenceNumber: generateReferenceNumber(),
		Fee:             0.0,
	}

	// Calculate balance after transaction
	switch tType {
	case Deposit:
		if toAcc != "" {
			if balance, err := ts.bankingSystem.accounts.GetBalance(toAcc); err == nil {
				transaction.BalanceAfter = balance + amount
			}
		}
	case Withdrawal:
		if fromAcc != "" {
			if balance, err := ts.bankingSystem.accounts.GetBalance(fromAcc); err == nil {
				transaction.BalanceAfter = balance - amount
			}
		}
	case Transfer:
		if fromAcc != "" {
			if balance, err := ts.bankingSystem.accounts.GetBalance(fromAcc); err == nil {
				transaction.BalanceAfter = balance - amount
			}
		}
	}

	transaction.Status = Completed
	ts.transactions[transaction.ID] = transaction

	return transaction, nil
}

func (ts *transactionService) GetTransaction(transactionID string) (*Transaction, error) {
	transaction, exists := ts.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}
	return transaction, nil
}

func (ts *transactionService) GetTransactionsByAccount(accountNumber string) ([]*Transaction, error) {
	var accountTransactions []*Transaction

	for _, transaction := range ts.transactions {
		if transaction.FromAccount == accountNumber || transaction.ToAccount == accountNumber {
			accountTransactions = append(accountTransactions, transaction)
		}
	}

	if len(accountTransactions) == 0 {
		return nil, errors.New("no transactions found for account")
	}

	return accountTransactions, nil
}

func (ts *transactionService) GetTransactionsByUser(userID int) ([]*Transaction, error) {
	user, err := ts.bankingSystem.users.Get(userID)
	if err != nil {
		return nil, err
	}

	var userTransactions []*Transaction

	for _, accountNumber := range user.Accounts {
		accountTxns, err := ts.GetTransactionsByAccount(accountNumber)
		if err == nil {
			userTransactions = append(userTransactions, accountTxns...)
		}
	}

	if len(userTransactions) == 0 {
		return nil, errors.New("no transactions found for user")
	}

	return userTransactions, nil
}

func (ts *transactionService) GetTransactionsByType(tType TransactionType) ([]*Transaction, error) {
	var typeTransactions []*Transaction

	for _, transaction := range ts.transactions {
		if transaction.Type == tType {
			typeTransactions = append(typeTransactions, transaction)
		}
	}

	if len(typeTransactions) == 0 {
		return nil, fmt.Errorf("no %s transactions found", tType)
	}

	return typeTransactions, nil
}

func (ts *transactionService) GetTransactionsByDateRange(startDate, endDate time.Time) ([]*Transaction, error) {
	var dateRangeTransactions []*Transaction

	for _, transaction := range ts.transactions {
		if (transaction.Timestamp.Equal(startDate) || transaction.Timestamp.After(startDate)) &&
			(transaction.Timestamp.Equal(endDate) || transaction.Timestamp.Before(endDate)) {
			dateRangeTransactions = append(dateRangeTransactions, transaction)
		}
	}

	if len(dateRangeTransactions) == 0 {
		return nil, errors.New("no transactions found in the specified date range")
	}

	return dateRangeTransactions, nil
}

func (ts *transactionService) UpdateTransactionStatus(transactionID string, status TransactionStatus) error {
	transaction, exists := ts.transactions[transactionID]
	if !exists {
		return errors.New("transaction not found")
	}

	transaction.Status = status
	ts.transactions[transactionID] = transaction
	return nil
}

func (ts *transactionService) GetAllTransactions() []*Transaction {
	allTransactions := make([]*Transaction, 0, len(ts.transactions))
	for _, transaction := range ts.transactions {
		allTransactions = append(allTransactions, transaction)
	}
	return allTransactions
}

func (ts *transactionService) GetTransactionSummary(accountNumber string) *TransactionSummary {
	summary := &TransactionSummary{
		AccountNumber: accountNumber,
	}

	transactions, err := ts.GetTransactionsByAccount(accountNumber)
	if err != nil {
		return summary
	}

	for _, transaction := range transactions {
		if transaction.Status == Completed {
			switch transaction.Type {
			case Deposit:
				summary.TotalDeposits += transaction.Amount
			case Withdrawal:
				summary.TotalWithdrawals += transaction.Amount
			case Transfer:
				if transaction.FromAccount == accountNumber {
					summary.TotalTransfersOut += transaction.Amount
				} else {
					summary.TotalTransfersIn += transaction.Amount
				}
			case Fee:
				summary.TotalFees += transaction.Amount
			}
			summary.TransactionCount++

			if transaction.Timestamp.After(summary.LastTransaction) {
				summary.LastTransaction = transaction.Timestamp
			}
		}
	}

	return summary
}

func (t Transaction) DisplayTransaction() {
	fmt.Println("=== Transaction Details ===")
	fmt.Printf("Transaction ID: %s\n", t.ID)
	fmt.Printf("Reference: %s\n", t.ReferenceNumber)
	fmt.Printf("Type: %s\n", t.Type)
	fmt.Printf("Status: %s\n", t.Status)

	if t.FromAccount != "" {
		fmt.Printf("From: %s\n", t.FromAccount)
	}
	if t.ToAccount != "" {
		fmt.Printf("To: %s\n", t.ToAccount)
	}

	fmt.Printf("Amount: $%.2f\n", t.Amount)

	if t.Fee > 0 {
		fmt.Printf("Fee: $%.2f\n", t.Fee)
	}

	if t.BalanceAfter > 0 {
		fmt.Printf("Balance After: $%.2f\n", t.BalanceAfter)
	}

	fmt.Printf("Time: %s\n", t.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Description: %s\n", t.Description)
	fmt.Println("---------------------------")
}

func (ts TransactionSummary) DisplayTransactionSummary() {
	fmt.Printf("\n=== Transaction Summary for Account: %s ===\n", ts.AccountNumber)
	fmt.Printf("Total Transactions: %d\n", ts.TransactionCount)
	fmt.Printf("Total Deposits: $%.2f\n", ts.TotalDeposits)
	fmt.Printf("Total Withdrawals: $%.2f\n", ts.TotalWithdrawals)
	fmt.Printf("Total Transfers Out: $%.2f\n", ts.TotalTransfersOut)
	fmt.Printf("Total Transfers In: $%.2f\n", ts.TotalTransfersIn)
	fmt.Printf("Total Fees: $%.2f\n", ts.TotalFees)

	if !ts.LastTransaction.IsZero() {
		fmt.Printf("Last Transaction: %s\n", ts.LastTransaction.Format("2006-01-02 15:04:05"))
	}

	netAmount := ts.TotalDeposits + ts.TotalTransfersIn - ts.TotalWithdrawals - ts.TotalTransfersOut - ts.TotalFees
	fmt.Printf("Net Amount: $%.2f\n", netAmount)
	fmt.Println("-----------------------------------")
}

func DisplayTransactions(transactions []*Transaction, title string) {
	fmt.Printf("\n=== %s ===\n", title)
	if len(transactions) == 0 {
		fmt.Println("No transactions found.")
		return
	}

	for _, transaction := range transactions {
		transaction.DisplayTransaction()
	}
}