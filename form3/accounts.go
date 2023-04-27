package form3

import (
	"context"
	"fmt"
)

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
	Version        *int64             `json:"version,omitempty"`
	CreatedOn      string             `json:"created_on,omitempty"`
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
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
	CustomerID              string   `json:"customer_id,omitempty"`
}

func (s *AccountService) Create(ctx context.Context, account *Account) (*Account, error) {
	request, error := s.client.NewRequest("POST", "/v1/organisation/accounts", account)

	if error != nil {
		return &Account{}, error
	}

	_, error = s.client.Do(ctx, request, account)

	if error != nil {
		return &Account{}, error
	}

	return account, error
}

func (s *AccountService) Fetch(ctx context.Context, accountId string) (*Account, error) {
	request, error := s.client.NewRequest("GET", fmt.Sprintf("/v1/organisation/accounts/%s", accountId), nil)

	if error != nil {
		return &Account{}, error
	}

	account := &Account{}
	_, error = s.client.Do(ctx, request, account)

	if error != nil {
		return &Account{}, error
	}

	return account, error
}
