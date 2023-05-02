package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

const resourceUri string = "/v1/organisation/accounts"

type AccountService struct {
	client *Client
}

type Account struct {
	Data *AccountData `json:"data,omitempty"`
}

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        int64              `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   string   `json:"account_classification,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 string   `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  string   `json:"status,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message,omitempty"`
}

func (s *AccountService) Create(account *Account) (*Account, error) {
	requestURL := fmt.Sprintf("%s%s", s.client.baseUrl, resourceUri)

	body, error := json.Marshal(&account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem marshalling the request body: %w", error)
	}

	response, error := s.client.httpClient.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if error != nil {
		return nil, fmt.Errorf("there was a problem performing the request: %w", error)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return handleUnsuccessfulOperation(response)
	}

	account = &Account{}
	error = json.NewDecoder(response.Body).Decode(&account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem unmarshalling the response body: %w", error)
	}

	return account, error
}

func (s *AccountService) Fetch(accountId string) (*Account, error) {
	requestURL := fmt.Sprintf("%s%s/%s", s.client.baseUrl, resourceUri, accountId)

	response, error := s.client.httpClient.Get(requestURL)

	if error != nil {
		return nil, fmt.Errorf("there was a problem performing the request: %w", error)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return handleUnsuccessfulOperation(response)
	}

	account := &Account{}
	error = json.NewDecoder(response.Body).Decode(&account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem unmarshalling the response body: %w", error)
	}

	return account, error
}

func (s *AccountService) Delete(accountId string, version int64) error {
	requestURL := fmt.Sprintf("%s%s/%s/?version=%d", s.client.baseUrl, resourceUri, accountId, version)

	request, error := http.NewRequest(http.MethodDelete, requestURL, nil)

	if error != nil {
		return fmt.Errorf("there was a problem creating the request: %w", error)
	}

	response, error := s.client.httpClient.Do(request)

	if error != nil {
		return fmt.Errorf("there was a problem performing the request: %w", error)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		_, error = handleUnsuccessfulOperation(response)

		return error
	}

	return nil
}

func handleUnsuccessfulOperation(response *http.Response) (*Account, error) {
	dump, error := httputil.DumpResponse(response, true)

	if error != nil {
		return nil, fmt.Errorf("it was not possible dump the response for an unsucessful operation: %w", error)
	}

	return nil, OperationError{
		Message:  "could not perform operation",
		Response: dump,
	}
}
