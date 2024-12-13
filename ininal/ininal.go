package ininal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.ininal.com/v3.0"
)

type UserDetails struct {
	Name                       string        `json:"name"`
	Surname                    string        `json:"surname"`
	Email                      string        `json:"email"`
	GsmNumber                  string        `json:"gsmNumber"`
	BirthDate                  string        `json:"birthDate"`
	Status                     string        `json:"status"`
	KycStatus                  string        `json:"kycStatus"`
	KycProcessStatus           string        `json:"kycProcessStatus"`
	TotalActiveCardBalance     float64       `json:"totalActiveCardBalance"`
	AvailableCashdrawAmount    float64       `json:"availableCashdrawAmount"`
	Education                  string        `json:"education"`
	Profession                 string        `json:"profession"`
	InternationalPassportNo    string        `json:"internationalPassportNo"`
	UserStatusText             string        `json:"userStatusText"`
	EmailVerified              bool          `json:"emailVerified"`
	EmailAllowed               bool          `json:"emailAllowed"`
	PhoneAllowed               bool          `json:"phoneAllowed"`
	SmsAllowed                 bool          `json:"smsAllowed"`
	LoadableLimit              float64       `json:"loadableLimit"`
	MonthlyLoadableLimit       float64       `json:"monthlyLoadableLimit"`
	CashWithdrawLimit          float64       `json:"cashWithdrawLimit"`
	ActiveWalletCampaignIdList []interface{} `json:"activeWalletCampaignIdList"`
}

type CardAccountResponse struct {
	Response    CardAccount `json:"response"`
	Description string      `json:"description"`
	HttpCode    int         `json:"httpCode"`
}

type CardAccount struct {
	LoadableLimit            float64       `json:"loadableLimit"`
	MonthlyLoadableLimit     float64       `json:"monthlyLoadableLimit"`
	ExchangeDailySellCount   int           `json:"exchangeDailySellCount"`
	ExchangeDailyBuyCount    int           `json:"exchangeDailyBuyCount"`
	ExchangeMonthlySellCount int           `json:"exchangeMonthlySellCount"`
	ExchangeMonthlyBuyCount  int           `json:"exchangeMonthlyBuyCount"`
	CashdrawBlockedAmount    float64       `json:"cashdrawBlockedAmount"`
	AvailableCashdrawAmount  float64       `json:"availableCashdrawAmount"`
	AccountListResponse      []AccountInfo `json:"accountListResponse"`
	AccessToken              string        `json:"accessToken`
}

type AccountInfo struct {
	AccountNumber  string  `json:"accountNumber"`
	AccountName    string  `json:"accountName"`
	AccountStatus  string  `json:"accountStatus"`
	AccountBalance float64 `json:"accountBalance"`
	IsFavorite     bool    `json:"isFavorite"`
	Currency       string  `json:"currency"`
	Iban           string  `json:"iban"`

	IbanValid        bool       `json:"ibanValid"`
	CardListResponse []CardInfo `json:"cardListResponse"`
	AvailableBalance float64    `json:"availableBalance"`
}

type CardInfo struct {
	CardId        int    `json:"cardId"`
	ProductCode   string `json:"productCode"`
	CardStatus    string `json:"cardStatus"`
	CardType      string `json:"cardType"`
	BarcodeNumber string `json:"barcodeNumber"`
	CardNumber    string `json:"cardNumber"`
	CardToken     string `json:"cardToken"`
}

type Transaction struct {
	TransactionDate  time.Time `json:"transactionDate"`
	Description      string    `json:"description"`
	ReferenceNo      string    `json:"referenceNo"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	Icon             string    `json:"icon"`
	TransactionType  string    `json:"transactionType"`
	RepeatActionType string    `json:"repeatActionType"`
}
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) GetUserDetails(userToken, authToken string) (*UserDetails, error) {
	url := fmt.Sprintf("https://api.ininal.com/v3.0/users/%s", userToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Host", "api.ininal.com")
	req.Header.Set("Content-Language", "en")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ininal/3.7.2 (com.ngier.ininalwallet; build:1; iOS 18.1.0) Alamofire/5.4.4")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Accept-Language", "en-US;q=1.0, ja-US;q=0.9, de-US;q=0.8")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Response UserDetails `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Response, nil
}

type LoginRequest struct {
	Password        string `json:"password"`
	DeviceSignature string `json:"deviceSignature"`
	DeviceID        string `json:"deviceId"`
	DeviceName      string `json:"deviceName"`
	LoginCredential string `json:"loginCredential"`
	AppVersion      string `json:"appVersion"`
	Token           string `json:"token"`
}

type LoginResponse struct {
	HTTPCode    int    `json:"httpCode"`
	Description string `json:"description"`
	Response    struct {
		AuthStatus string `json:"authStatus"`
		Token      string `json:"token"`
		UserToken  string `json:"userToken"`
	} `json:"response"`
	ValidationErrors interface{} `json:"validationErrors"`
}

type VerifyRequest struct {
	OTP   string `json:"otp"`
	Token string `json:"token"`
}

func (c *Client) Login(password, deviceID, loginCredential, loginToken string, bearerToken string, deviceSignature string) (*LoginResponse, error) {

	req := LoginRequest{
		Password:        password,
		DeviceSignature: deviceSignature,
		DeviceID:        deviceID,
		DeviceName:      "iPhone16,1",
		LoginCredential: loginCredential,
		AppVersion:      "3.7.6",
		Token:           loginToken,
	}

	fmt.Println("logging in with")
	fmt.Println("password:", password)
	fmt.Println("deviceID:", deviceID)
	fmt.Println("loginCredential:", loginCredential)
	fmt.Println("token:", loginToken)
	fmt.Println("bearer:", bearerToken)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	request, err := http.NewRequest("POST", BaseURL+"/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+bearerToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Language", "en")
	request.Header.Set("User-Agent", "ininal/3.7.6 (com.ngier.ininalwallet; build:2; iOS 18.2.0) Alamofire/5.4.4")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &loginResp, nil
}

func (c *Client) Verify(otp, token string, bearerToken string) (*LoginResponse, error) {
	req := VerifyRequest{
		OTP:   otp,
		Token: token,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	request, err := http.NewRequest("POST", BaseURL+"/auth/login/verify", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+bearerToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Language", "en")
	request.Header.Set("User-Agent", "ininal/3.7.6 (com.ngier.ininalwallet; build:2; iOS 18.2.0) Alamofire/5.4.4")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var verifyResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &verifyResp, nil
}

// generateDeviceSignature creates an RSA signature for the device ID
// Note: You'll need to implement the actual signature generation logic based on the app's requirements
func (c *Client) generateDeviceSignature(deviceID string) (string, error) {
	// Implementation depends on the actual signature algorithm used by the app
	// You'll need to reverse engineer this from the app
	return "", nil
}

func (c *Client) GetUserCardAccount(deviceID, userToken, authToken string) (*CardAccount, error) {
	url := fmt.Sprintf("https://api.ininal.com/v3.2/users/%s/cardaccount", userToken)

	fmt.Println("Getting card account for user token:", userToken)
	fmt.Println("Auth token:", authToken)
	fmt.Println("URL:", url)
	fmt.Println("Device ID:", deviceID)

	reqBody, err := json.Marshal(map[string]string{
		"deviceId": deviceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Host", "api.ininal.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Language", "en")
	req.Header.Set("Accept-Language", "en-US;q=1.0, ja-US;q=0.9, de-US;q=0.8")
	req.Header.Set("User-Agent", "ininal/3.7.2 (com.ngier.ininalwallet; build:1; iOS 18.1.0) Alamofire/5.4.4")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// copy body into a separate buffer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	// create a new buffer with the copied body
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// print body
	fmt.Println(string(body))

	var result CardAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result.Response, nil
}

func (c *Client) GetUserTransactions(userToken, authToken, accountID string, startDate, endDate time.Time, resultLimit int) ([]Transaction, error) {
	url := fmt.Sprintf("https://api.ininal.com/v3.1/users/%s/transactions/%s", userToken, accountID)

	fmt.Println(url)
	fmt.Println(authToken)

	if resultLimit == 0 {
		resultLimit = 3
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"startDate":   startDate.Format("2006/01/02"),
		"endDate":     endDate.Format("2006/01/02"),
		"resultLimit": resultLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Language", "en")
	req.Header.Set("Accept-Language", "en-US;q=1.0, ja-US;q=0.9, de-US;q=0.8")
	req.Header.Set("User-Agent", "ininal/3.7.2 (com.ngier.ininalwallet; build:1; iOS 18.1.0) Alamofire/5.4.4")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// b := bytes.Buffer{}
	// b.ReadFrom(resp.Body)
	// fmt.Println(b.String())
	// resp.Body = io.NopCloser(bytes.NewBuffer(b.Bytes()))

	var result struct {
		Response struct {
			TransactionList []Transaction `json:"transactionList"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Response.TransactionList, nil
}

type CustomerDetails struct {
	CustomerID                    int64       `json:"customerId"`
	CustomerToken                 string      `json:"customerToken"`
	Name                          string      `json:"name"`
	Surname                       string      `json:"surname"`
	Email                         string      `json:"email"`
	TCIdentificationNumber        interface{} `json:"tcIdentificationNumber"`
	GsmNumber                     string      `json:"gsmNumber"`
	BirthDate                     string      `json:"birthDate"`
	Password                      string      `json:"password"`
	Status                        string      `json:"status"`
	MotherMaidenName              string      `json:"motherMaidenName"`
	RegistrationChannel           string      `json:"registrationChannel"`
	LoadableLimit                 float64     `json:"loadableLimit"`
	MonthlyLoadableLimit          float64     `json:"monthlyLoadableLimit"`
	CashWithdrawLimit             float64     `json:"cashWithdrawLimit"`
	LoadableLimitDefault          float64     `json:"loadableLimitDefault"`
	MonthlyLoadableLimitDefault   float64     `json:"monthlyLoadableLimitDefault"`
	CashWithdrawLimitDefault      float64     `json:"cashWithdrawLimitDefault"`
	MaxAssignCardLimit            int         `json:"maxAssignCardLimit"`
	MaxAssignCardLimitDefault     int         `json:"maxAssignCardLimitDefault"`
	MaxActiveAccountsLimit        int         `json:"maxActiveAccountsLimit"`
	MaxActiveAccountsLimitDefault int         `json:"maxActiveAccountsLimitDefault"`
	KycStatus                     string      `json:"kycStatus"`
	KycProcessStatus              string      `json:"kycProcessStatus"`
	ShowMenu                      bool        `json:"showMenu"`
	TotalActiveCardBalance        float64     `json:"totalActiveCardBalance"`
	AvailableCashdrawAmount       float64     `json:"availableCashdrawAmount"`
	CashdrawBlockedAmount         float64     `json:"cashdrawBlockedAmount"`
	TempDate                      interface{} `json:"tempDate"`
	ActiveWalletCampaignIdList    []int       `json:"activeWalletCampaignIdList"`
	MoneyRequestDisabled          bool        `json:"moneyRequestDisabled"`
	Education                     string      `json:"education"`
	Profession                    string      `json:"profession"`
	ChildApprovalPending          bool        `json:"childApprovalPending"`
	ParentFullName                interface{} `json:"parentFullName"`
	ParentGsm                     interface{} `json:"parentGsm"`
	InternationalPassportNo       string      `json:"internationalPassportNo"`
	BeforeStatus                  interface{} `json:"beforeStatus"`
	RegistrationDate              interface{} `json:"registrationDate"`
	UserStatusText                string      `json:"userStatusText"`
	KvkkChecked                   bool        `json:"kvkkChecked"`
	SmsAllowed                    bool        `json:"smsAllowed"`
	EmailAllowed                  bool        `json:"emailAllowed"`
	PhoneAllowed                  bool        `json:"phoneAllowed"`
	EmailVerified                 bool        `json:"emailVerified"`
	UserAgreementApproved         bool        `json:"userAgreementApproved"`
	MotherMaidenNameEmpty         bool        `json:"motherMaidenNameEmpty"`
}

func (c *Client) GetCustomerDetails(userToken, authToken string) (*CustomerDetails, error) {
	url := fmt.Sprintf("https://api.ininal.com/v3.0/users/%s", userToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Language", "en")
	req.Header.Set("User-Agent", "ininal/3.7.2 (com.ngier.ininalwallet; build:1; iOS 18.1.0) Alamofire/5.4.4")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// copy body and print it out
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println(string(body))

	// put the body back into the response body
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	var result struct {
		HTTPCode    int             `json:"httpCode"`
		Description string          `json:"description"`
		Response    CustomerDetails `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result.Response, nil
}
