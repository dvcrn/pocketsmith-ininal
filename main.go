package main

import (
	"bytes"
	"encoding/json"
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
	AnapayUsername   string
	AnapayPassword   string
	PocketsmithToken string

	NumTransactions int
}

func getConfig() *Config {
	config := &Config{}

	// Define command-line flags
	// flag.StringVar(&config.AnapayUsername, "username", os.Getenv("ININAL_USERNAME"), "Ininal username")
	// flag.StringVar(&config.AnapayPassword, "password", os.Getenv("ININAL_PASSWORD"), "Ininal password")
	flag.StringVar(&config.PocketsmithToken, "token", os.Getenv("POCKETSMITH_TOKEN"), "Pocketsmith API token")
	// flag.IntVar(&config.NumTransactions, "num-transactions", 100, "Number of transactions to parse")
	flag.Parse()

	// Validate required fields
	// if config.AnapayUsername == "" {
	// 	fmt.Println("Error: Anapay username is required. Set via -username flag or ANAPAY_USERNAME environment variable")
	// 	os.Exit(1)
	// }
	// if config.AnapayPassword == "" {
	// 	fmt.Println("Error: Anapay password is required. Set via -password flag or ANAPAY_PASSWORD environment variable")
	// 	os.Exit(1)
	// }
	if config.PocketsmithToken == "" {
		fmt.Println("Error: Pocketsmith token is required. Set via -token flag or POCKETSMITH_TOKEN environment variable")
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
		"244078",
		"4FE92B1D-9D75-47C0-BD65-C650F8921441",
		"+818099784527",
		LOGIN_TOKEN, // login token??
		// USER_TOKEN,
		LOGIN_BEARER_TOKEN,
	)

	if err != nil {
		panic(err)
	}

	// TODO: for debug. remove me.
	func(v interface{}) {
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		buf := bytes.NewBuffer(j)
		fmt.Printf("%v\n", buf.String())
	}(loginResp)

	// wait and ask for OTP
	if loginResp.Response.AuthStatus == "OTP_REQUIRED" {
		fmt.Println("Please enter the OTP code sent to your phone:")
		var otp string
		fmt.Scanln(&otp)

		verifyResp, err := client.Verify(otp, loginResp.Response.Token, LOGIN_BEARER_TOKEN)
		if err != nil {
			panic(err)
		}

		// TODO: for debug. remove me.
		func(v interface{}) {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			buf := bytes.NewBuffer(j)
			fmt.Printf("%v\n", buf.String())
		}(verifyResp)

		userToken = verifyResp.Response.UserToken
		userAuth = verifyResp.Response.Token
	} else {
		fmt.Println("OTP not required")
		userToken = loginResp.Response.UserToken
		userAuth = loginResp.Response.Token
	}

	// GetDetails:
	fmt.Println("Got:")
	fmt.Println("userToken: ", userToken)
	fmt.Println("userAuth: ", userAuth)

	userCardAccount, err := client.GetUserDetails(userToken, userAuth)
	if err != nil {
		panic(err)
	}

	// TODO: for debug. remove me.
	func(v interface{}) {
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		buf := bytes.NewBuffer(j)
		fmt.Printf("%v\n", buf.String())
	}(userCardAccount)

	cardAccount, err := client.GetUserCardAccount(deviceID, userToken, userAuth)
	if err != nil {
		panic(err)
	}

	// Debug print card account details
	func(v interface{}) {
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		buf := bytes.NewBuffer(j)
		fmt.Printf("%v\n", buf.String())
	}(cardAccount)

	// Get transactions for each account
	for _, account := range cardAccount.AccountListResponse {
		fmt.Println("Creating Pocketsmith account for account", account.AccountNumber)

		psAcc, err := findOrCreateAccount(ps, res.ID, account.AccountName)
		if err != nil {
			fmt.Printf("Error creating/finding Pocketsmith account: %v\n", err)
			continue
		}
		fmt.Printf("Updating Ininal Account Balance to %f.2 for date %s\n", account.AccountBalance, time.Now().Format("2006-01-02"))

		updateRes, err := ps.UpdateTransactionAccount(psAcc.ID, psAcc.PrimaryTransactionAccount.Institution.ID, account.AccountBalance, time.Now().Format("2006-01-02"))
		if err != nil {
			fmt.Printf("Error updating Ininal account balance: %v\n", err)
			continue
		}
		fmt.Printf("Updated Ininal account balance: %.2f\n", updateRes.CurrentBalance)

		fmt.Printf("\nFetching transactions for account: %s\n", account.AccountName)

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
		for _, transaction := range transactions {
			fmt.Printf("Transaction: %s (%s) from %s\n", transaction.Description, transaction.ReferenceNo, transaction.TransactionDate.Format("2006-01-02"))
			if repeatedExistingTransactions > 10 {
				fmt.Println("Too many repeated existing transactions, exiting")
				break
			}

			startDate := transaction.TransactionDate.Add(-2 * 24 * time.Hour).Format("2006-01-02")
			endDate := transaction.TransactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")
			fmt.Printf("Searching for transactions between %s and %s\n", startDate, endDate)

			searchRes, err := ps.SearchTransactions(psAcc.PrimaryTransactionAccount.ID, startDate, endDate, "")
			if err != nil {
				fmt.Printf("Error searching for transaction: %v\n", err)
				continue
			}

			if len(searchRes) > 0 {
				for _, tx := range searchRes {
					checkNum := ""
					if tx.ChequeNumber != nil {
						checkNum = *tx.ChequeNumber
					}
					memo := ""
					if tx.Memo != nil {
						memo = *tx.Memo
					}

					if checkNum == transaction.ReferenceNo || memo == transaction.ReferenceNo {
						fmt.Println("Found transaction already, won't add it again: ", transaction.ReferenceNo)
						repeatedExistingTransactions++
						break
					}
				}
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
