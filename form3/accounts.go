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
	Client        *Client       // Used to access basic request configurations and perform http requests.
	JsonMarshal   JsonMarshal   // Used to marshal json.
	JsonUnmarshal JsonUnmarshal // Used to unmarshal json.
	ReadAll       ReadAll       // Used to read the response body of a http request.
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
func (s *AccountService) Create(account *Account) (*Account, *http.Response, error) {
	requestURL := fmt.Sprintf("%s%s", s.Client.BaseUrl, resourceUri)

	body, error := s.JsonMarshal(account)

	if error != nil {
		return nil, nil, OperationError{Message: error.Error()}
	}

	return s.handleAccountResponse(http.MethodPost, requestURL, body, http.StatusCreated)
}

// Create allows one to fetch a FORM3 account.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/fetch-an-account
func (s *AccountService) Fetch(accountId string) (*Account, *http.Response, error) {
	requestURL := fmt.Sprintf("%s%s/%s", s.Client.BaseUrl, resourceUri, accountId)

	return s.handleAccountResponse(http.MethodGet, requestURL, nil, http.StatusOK)
}

// Create allows one to delete a FORM3 account.
//
// More details available in: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/delete-an-account
func (s *AccountService) Delete(accountId string, version int64) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s%s/%s?version=%d", s.Client.BaseUrl, resourceUri, accountId, version)

	response, error := s.Client.PerformRequest(http.MethodDelete, requestURL, nil)

	if error != nil {
		return nil, error
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		body, error := s.ReadAll(response.Body)

		if error != nil {
			return response, OperationError{Message: error.Error()}
		}

		return response, OperationError{
			Message: response.Status,
			Body:    body,
		}
	}

	return response, nil
}

func (s *AccountService) handleAccountResponse(httpMethod string, requestURL string, body []byte, successfulStatusCode int) (*Account, *http.Response, error) {
	response, error := s.Client.PerformRequest(httpMethod, requestURL, body)

	if error != nil {
		return nil, nil, error
	}

	defer response.Body.Close()

	body, error = s.ReadAll(response.Body)

	if error != nil {
		return nil, response, OperationError{Message: error.Error()}
	}

	if response.StatusCode != successfulStatusCode {
		return nil, response, OperationError{
			Message: response.Status,
			Body:    body,
		}
	}

	account := &Account{}
	error = s.JsonUnmarshal(body, &account)

	if error != nil {
		return nil, response, OperationError{Message: error.Error()}
	}

	return account, response, error
}
