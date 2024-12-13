package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dvcrn/pocketsmith-anapay/ininal"
	"github.com/dvcrn/pocketsmith-go"
)

const INSTITUION_NAME = "Ininal"
const ACCOUNT_NAME = "Ininal"

type Config struct {
	DeviceID         string
	LoginToken       string
	UserToken        string
	LoginBearerToken string
	PocketsmithToken string

	Password        string
	LoginCredential string
	DeviceSignature string
}

func getConfig() *Config {
	config := &Config{}

	// Define command-line flags
	flag.StringVar(&config.DeviceID, "device-id", os.Getenv("ININAL_DEVICE_ID"), "Ininal device ID")
	flag.StringVar(&config.LoginToken, "login-token", os.Getenv("ININAL_LOGIN_TOKEN"), "Ininal login token")
	flag.StringVar(&config.UserToken, "user-token", os.Getenv("ININAL_USER_TOKEN"), "Ininal user token")
	flag.StringVar(&config.LoginBearerToken, "login-bearer-token", os.Getenv("ININAL_LOGIN_BEARER_TOKEN"), "Ininal login bearer token")

	flag.StringVar(&config.PocketsmithToken, "pocketsmith-token", os.Getenv("POCKETSMITH_TOKEN"), "Pocketsmith API token")

	flag.StringVar(&config.Password, "password", os.Getenv("ININAL_PASSWORD"), "Ininal password (App PIN)")
	flag.StringVar(&config.LoginCredential, "login-credential", os.Getenv("ININAL_LOGIN_CREDENTIAL"), "Ininal login credential (Phone Number)")
	flag.StringVar(&config.DeviceSignature, "device-signature", os.Getenv("ININAL_DEVICE_SIGNATURE"), "Ininal device signature")

	flag.Parse()

	// Validate required fields
	if config.DeviceID == "" {
		fmt.Println("Error: Device ID is required. Set via -device-id flag or ININAL_DEVICE_ID environment variable")
		os.Exit(1)
	}
	if config.LoginToken == "" {
		fmt.Println("Error: Login token is required. Set via -login-token flag or ININAL_LOGIN_TOKEN environment variable")
		os.Exit(1)
	}
	if config.UserToken == "" {
		fmt.Println("Error: User token is required. Set via -user-token flag or ININAL_USER_TOKEN environment variable")
		os.Exit(1)
	}
	if config.LoginBearerToken == "" {
		fmt.Println("Error: Login bearer token is required. Set via -login-bearer-token flag or ININAL_LOGIN_BEARER_TOKEN environment variable")
		os.Exit(1)
	}
	if config.PocketsmithToken == "" {
		fmt.Println("Error: Pocketsmith token is required. Set via -token flag or POCKETSMITH_TOKEN environment variable")
		os.Exit(1)
	}

	if config.Password == "" {
		fmt.Println("Error: Password is required. Set via -password flag or ININAL_PASSWORD environment variable")
		os.Exit(1)
	}
	if config.LoginCredential == "" {
		fmt.Println("Error: Login credential is required. Set via -login-credential flag or ININAL_LOGIN_CREDENTIAL environment variable")
		os.Exit(1)
	}
	if config.DeviceSignature == "" {
		fmt.Println("Error: Device signature is required. Set via -device-signature flag or ININAL_DEVICE_SIGNATURE environment variable")
		os.Exit(1)
	}

	return config
}

func findOrCreateAccount(ps *pocketsmith.Client, userID int, accountName string) (*pocketsmith.Account, error) {
	name := fmt.Sprintf("Ininal %s", accountName)
	account, err := ps.FindAccountByName(userID, name)
	if err != nil {
		if err != pocketsmith.ErrNotFound {
			return nil, err
		}

		institution, err := ps.FindInstitutionByName(userID, INSTITUION_NAME)
		if err != nil {
			if err != pocketsmith.ErrNotFound {
				return nil, err
			}

			institution, err = ps.CreateInstitution(userID, INSTITUION_NAME, "try")
			if err != nil {
				return nil, err
			}
		}

		account, err := ps.CreateAccount(userID, institution.ID, name, "try", pocketsmith.AccountTypeCredits)
		if err != nil {
			return nil, err
		}

		return account, nil
	}

	return account, nil
}

func main() {
	config := getConfig()

	ps := pocketsmith.NewClient(config.PocketsmithToken)
	res, err := ps.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	fmt.Println("Pocketsmith user ID:", res.ID)

	client := ininal.NewClient()

	var userToken string
	var userAuth string

	// First login step
	loginResp, err := client.Login(
		config.Password,
		config.DeviceID,
		config.LoginCredential,
		config.LoginToken,
		config.LoginBearerToken,
		config.DeviceSignature,
	)

	if err != nil {
		panic(err)
	}

	// wait and ask for OTP
	if loginResp.Response.AuthStatus == "OTP_REQUIRED" {
		fmt.Println("Please enter the OTP code sent to your phone:")
		var otp string
		fmt.Scanln(&otp)

		verifyResp, err := client.Verify(otp, loginResp.Response.Token, config.LoginBearerToken)
		if err != nil {
			panic(err)
		}

		userToken = verifyResp.Response.UserToken
		userAuth = verifyResp.Response.Token
	} else {
		fmt.Println("OTP not required")
		userToken = loginResp.Response.UserToken
		userAuth = loginResp.Response.Token
	}

	// GetDetails:
	// fmt.Println("Got:")
	// fmt.Println("userToken: ", userToken)
	// fmt.Println("userAuth: ", userAuth)

	cardAccount, err := client.GetUserCardAccount(config.DeviceID, userToken, userAuth)
	if err != nil {
		panic(err)
	}

	// Get transactions for each account
	for _, account := range cardAccount.AccountListResponse {
		fmt.Println("Creating Pocketsmith account for account", account.AccountNumber)

		psAcc, err := findOrCreateAccount(ps, res.ID, account.AccountName)
		if err != nil {
			fmt.Printf("Error creating/finding Pocketsmith account: %v\n", err)
			continue
		}

		dateString := time.Now().Format("2006-01-02")

		updateRes, err := ps.UpdateTransactionAccount(psAcc.PrimaryTransactionAccount.ID, psAcc.PrimaryTransactionAccount.Institution.ID, account.AccountBalance, dateString)
		if err != nil {
			fmt.Printf("Error updating Ininal account balance: %v\n", err)
			continue
		}
		fmt.Println("Updated Ininal Account balance: ", updateRes.CurrentBalance)

		transactions, err := client.GetUserTransactions(
			userToken,
			cardAccount.AccessToken,
			account.AccountNumber,
			time.Now().AddDate(-2, 0, 0), // Start date
			time.Now(),
			200, // Number of transactions
		)
		if err != nil {
			fmt.Printf("Error fetching transactions: %v\n", err)
			continue
		}

		fmt.Println(len(transactions))

		repeatedExistingTransactions := 0
		for i, transaction := range transactions {
			fmt.Printf("[%d/%d] Transaction: %s (%s) from %s\n", i+1, len(transactions), transaction.Description, transaction.ReferenceNo, transaction.TransactionDate.Format("2006-01-02"))

			if repeatedExistingTransactions > 10 {
				fmt.Println("Too many repeated existing transactions, exiting")
				break
			}

			searchRes, err := ps.SearchTransactionsByMemoContains(psAcc.PrimaryTransactionAccount.ID, transaction.TransactionDate, transaction.ReferenceNo)
			if err != nil {
				fmt.Printf("Error searching for existing transaction: %v\n", err)
				continue
			}

			if len(searchRes) > 0 {
				fmt.Println("Found existing transaction by ref number: ", transaction.ReferenceNo)
				repeatedExistingTransactions++
				continue
			}

			createTx := &pocketsmith.CreateTransaction{
				Payee:        strings.TrimSpace(transaction.Description),
				Amount:       transaction.Amount,
				Date:         transaction.TransactionDate.Format("2006-01-02"),
				IsTransfer:   strings.Contains(transaction.TransactionType, "Banka Transferi"),
				ChequeNumber: transaction.ReferenceNo,
				Note:         transaction.TransactionType,
				Memo:         transaction.ReferenceNo,
			}

			fmt.Println("Creating transaction with createTx: ", createTx.Payee, createTx.Amount, createTx.Date, createTx.IsTransfer, createTx.Note)
			_, err = ps.AddTransaction(psAcc.PrimaryTransactionAccount.ID, createTx)
			if err != nil {
				fmt.Printf("Error creating transaction: %v\n", err)
				continue
			}
		}
	}
}
