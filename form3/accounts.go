package form3

import (
	"fmt"
	"io"
	"net/http"
)

// resourseUri contains the path to the resource.
const resourceUri string = "/v1/organisation/accounts"

// JsonMarshal defines the function interface that is used to marshal json.
type JsonMarshal func(v any) ([]byte, error)

// JsonUnmarshal defines the function interface that is used to unmarshal json.
type JsonUnmarshal func(data []byte, v any) error

// ReadAll defines the function interface that is used to read a response body.
type ReadAll func(r io.Reader) ([]byte, error)

// AccountService allows access to operations related to accounts.
type AccountService struct {
	Client        *Client
	JsonMarshal   JsonMarshal
	JsonUnmarshal JsonUnmarshal
	ReadAll       ReadAll
}

// Represents a FORM3 account.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts
type Account struct {
	Data *AccountData `json:"data,omitempty"`
}

// Represents a FORM3 account data.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        int64              `json:"version,omitempty"`
}

// Represents a FORM3 account attributes.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts
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

// Create allows one to create a FORM3 account.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/create-an-account
func (s *AccountService) Create(account *Account) (*Account, error) {
	requestURL := fmt.Sprintf("%s%s", s.Client.BaseUrl, resourceUri)

	body, error := s.marshal(account)

	if error != nil {
		return nil, error
	}

	response, error := PerformRequest(s.Client, http.MethodPost, requestURL, body)

	if error != nil {
		return nil, error
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, BuildUnsuccessfulResponse(response)
	}

	return s.unmarshal(response)
}

// Create allows one to fetch a FORM3 account.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/fetch-an-account
func (s *AccountService) Fetch(accountId string) (*Account, error) {
	requestURL := fmt.Sprintf("%s%s/%s", s.Client.BaseUrl, resourceUri, accountId)

	response, error := PerformRequest(s.Client, http.MethodGet, requestURL, nil)

	if error != nil {
		return nil, error
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, BuildUnsuccessfulResponse(response)
	}

	return s.unmarshal(response)
}

// Create allows one to delete a FORM3 account.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/delete-an-account
func (s *AccountService) Delete(accountId string, version int64) error {
	requestURL := fmt.Sprintf("%s%s/%s/?version=%d", s.Client.BaseUrl, resourceUri, accountId, version)

	response, error := PerformRequest(s.Client, http.MethodDelete, requestURL, nil)

	if error != nil {
		return fmt.Errorf("there was a problem performing the request: %w", error)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return BuildUnsuccessfulResponse(response)
	}

	return nil
}

func (s *AccountService) marshal(account *Account) ([]byte, error) {
	body, error := s.JsonMarshal(&account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem marshalling the request body: %w", error)
	}

	return body, nil
}

func (s *AccountService) unmarshal(response *http.Response) (*Account, error) {
	account := &Account{}
	data, error := s.ReadAll(response.Body)

	if error != nil {
		return nil, fmt.Errorf("there was a problem reading the response body: %w", error)
	}

	error = s.JsonUnmarshal(data, &account)

	if error != nil {
		return nil, fmt.Errorf("there was a problem unmarshalling the response body: %w", error)
	}

	return account, error
}
