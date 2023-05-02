package form3

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (s *AccountService) Create(account *Account) (*Account, error) {
	requestURL := fmt.Sprintf("%s%s", s.client.baseUrl, resourceUri)

	body, error := marshal(account)

	if error != nil {
		return nil, error
	}

	response, error := performRequest(s.client, http.MethodPost, requestURL, body)

	if error != nil {
		return nil, error
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, handleUnsuccessfulOperation(response)
	}

	return unmarshal(response)
}

func (s *AccountService) Fetch(accountId string) (*Account, error) {
	requestURL := fmt.Sprintf("%s%s/%s", s.client.baseUrl, resourceUri, accountId)

	response, error := performRequest(s.client, http.MethodGet, requestURL, nil)

	if error != nil {
		return nil, error
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, handleUnsuccessfulOperation(response)
	}

	return unmarshal(response)
}

func (s *AccountService) Delete(accountId string, version int64) error {
	requestURL := fmt.Sprintf("%s%s/%s/?version=%d", s.client.baseUrl, resourceUri, accountId, version)

	response, error := performRequest(s.client, http.MethodDelete, requestURL, nil)

	if error != nil {
		return fmt.Errorf("there was a problem performing the request: %w", error)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return handleUnsuccessfulOperation(response)
	}

	return nil
}

func marshal(account *Account) ([]byte, error) {
	body, error := json.Marshal(&account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem marshalling the request body: %w", error)
	}

	return body, nil
}

func unmarshal(response *http.Response) (*Account, error) {
	account := &Account{}
	error := json.NewDecoder(response.Body).Decode(&account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem unmarshalling the response body: %w", error)
	}

	return account, error
}
