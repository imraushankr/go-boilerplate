package main

import (
	"bank-system/bank"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	bankingSystem := bank.NewBankingSystem()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Integrated Banking System ===")
	fmt.Println("Welcome to the Banking System!")

	// createSampleData(bankingSystem)

	for {
		displayMenu()
		fmt.Print("Enter your choice: ")
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			createUserHandler(bankingSystem, scanner)
		case "2":
			createAccountHandler(bankingSystem, scanner)
		case "3":
			depositHandler(bankingSystem, scanner)
		case "4":
			withdrawHandler(bankingSystem, scanner)
		case "5":
			transferHandler(bankingSystem, scanner)
		case "6":
			viewAccountHandler(bankingSystem, scanner)
		case "7":
			viewUserHandler(bankingSystem, scanner)
		case "8":
			bankingSystem.ListAllAccounts()
		case "9":
			bankingSystem.ListAllUsers()
		case "10":
			bankingSystem.GetTransactionHistory()
		case "11":
			viewAccountTransactionsHandler(bankingSystem, scanner)
		case "12":
			viewTransactionSummaryHandler(bankingSystem, scanner)
		case "13":
			viewBalanceHandler(bankingSystem, scanner)
		case "14":
			closeAccountHandler(bankingSystem, scanner)
		case "15":
			fmt.Println("Exiting the Banking System. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func displayMenu() {
	fmt.Println("\nPlease select an option:")
	fmt.Println("1. Create User")
	fmt.Println("2. Create Account")
	fmt.Println("3. Deposit Money")
	fmt.Println("4. Withdraw Money")
	fmt.Println("5. Transfer Money")
	fmt.Println("6. View Account Details")
	fmt.Println("7. View User Details")
	fmt.Println("8. List All Accounts")
	fmt.Println("9. List All Users")
	fmt.Println("10. View All Transactions")
	fmt.Println("11. View Account Transactions")
	fmt.Println("12. View Transaction Summary")
	fmt.Println("13. View Account Balance")
	fmt.Println("14. Close Account")
	fmt.Println("15. Exit")
}

// func createSampleData(bs *bank.BankingSystem) {
// 	// Create sample user
// 	user, err := bs.CreateUser(
// 		"Raushan",
// 		"Kumar",
// 		"raushan.kumar@hk.com",
// 		"RK@tr",
// 		"Motihari, Bihar",
// 		"7645927364",
// 		"DFFGD7657JKHG",
// 		"3245733353443626",
// 	)
// 	if err != nil {
// 		fmt.Printf("Failed to create sample user: %v\n", err)
// 		return
// 	}

// 	// Create sample accounts for the user
// 	err = bs.CreateAccount("ACC001", "Raushan Kumar", "Savings", user.ID)
// 	if err != nil {
// 		fmt.Printf("Failed to create sample account: %v\n", err)
// 	}

// 	// Create some transactions
// 	bs.Deposit("ACC001", 1000.0)
// 	bs.Deposit("ACC001", 500.0)
// 	bs.Withdraw("ACC001", 200.0)
// }

func createUserHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Println("\n=== Create New User ===")

	fmt.Print("Enter first name: ")
	scanner.Scan()
	firstName := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter last name: ")
	scanner.Scan()
	lastName := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter email: ")
	scanner.Scan()
	email := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter address: ")
	scanner.Scan()
	address := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter phone: ")
	scanner.Scan()
	phone := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter PAN card number: ")
	scanner.Scan()
	panCard := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Aadhar card number: ")
	scanner.Scan()
	aadharCard := strings.TrimSpace(scanner.Text())

	_, err := bs.CreateUser(firstName, lastName, email, password, address, phone, panCard, aadharCard)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	}
}

func createAccountHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Println("\n=== Create New Account ===")

	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := strings.TrimSpace(scanner.Text())

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Println("Invalid user ID. Please enter a valid number.")
		return
	}

	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter holder name: ")
	scanner.Scan()
	holderName := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter account type (Savings/Current): ")
	scanner.Scan()
	accountType := strings.TrimSpace(scanner.Text())

	err = bs.CreateAccount(accountNumber, holderName, accountType, userID)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
	}
}

func depositHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter amount to deposit: ")
	scanner.Scan()
	amountStr := strings.TrimSpace(scanner.Text())

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Invalid amount. Please enter a valid number.")
		return
	}

	err = bs.Deposit(accountNumber, amount)
	if err != nil {
		fmt.Printf("Error depositing money: %v\n", err)
	}
}

func withdrawHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter amount to withdraw: ")
	scanner.Scan()
	amountStr := strings.TrimSpace(scanner.Text())

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Invalid amount. Please enter a valid number.")
		return
	}

	err = bs.Withdraw(accountNumber, amount)
	if err != nil {
		fmt.Printf("Error withdrawing money: %v\n", err)
	}
}

func transferHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter source account number: ")
	scanner.Scan()
	fromAccount := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter destination account number: ")
	scanner.Scan()
	toAccount := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter amount to transfer: ")
	scanner.Scan()
	amountStr := strings.TrimSpace(scanner.Text())

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Invalid amount. Please enter a valid number.")
		return
	}

	err = bs.Transfer(fromAccount, toAccount, amount)
	if err != nil {
		fmt.Printf("Error transferring money: %v\n", err)
	}
}

func viewAccountHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	account, err := bs.GetAccount(accountNumber)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\n=== Account Details ===")
	account.DisplayAccountInfo()
}

func viewUserHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := strings.TrimSpace(scanner.Text())

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Println("Invalid user ID. Please enter a valid number.")
		return
	}

	user, err := bs.GetUser(userID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	user.DisplayUserInfo()
}

func viewAccountTransactionsHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	bs.GetAccountTransactions(accountNumber)
}

func viewTransactionSummaryHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	bs.GetTransactionSummary(accountNumber)
}

func viewBalanceHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	balance, err := bs.GetBalance(accountNumber)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Account %s balance: $%.2f\n", accountNumber, balance)
}

func closeAccountHandler(bs *bank.BankingSystem, scanner *bufio.Scanner) {
	fmt.Print("Enter account number to close: ")
	scanner.Scan()
	accountNumber := strings.TrimSpace(scanner.Text())

	err := bs.CloseAccount(accountNumber)
	if err != nil {
		fmt.Printf("Error closing account: %v\n", err)
	} else {
		fmt.Printf("Account %s closed successfully\n", accountNumber)
	}
}