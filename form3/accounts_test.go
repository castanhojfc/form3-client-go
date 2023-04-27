package form3

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccounts_Create(t *testing.T) {
	t.Run("creates account with valid data", func(*testing.T) {
		client := NewClient()
		context := context.Background()

		account := &Account{
			Data: &AccountData{
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Type:           "accounts",
				Attributes: &AccountAttributes{
					Country:                 "GB",
					BaseCurrency:            "GBP",
					BankID:                  "400302",
					BankIDCode:              "GBDSC",
					CustomerID:              "234",
					Bic:                     "NWBKGB42",
					Name:                    []string{"Samantha Holder"},
					AlternativeNames:        []string{"Sam Holder"},
					AccountClassification:   "Personal",
					JointAccount:            false,
					AccountMatchingOptOut:   false,
					SecondaryIdentification: "A1B2C3D4",
				},
			},
		}

		account, error := client.Accounts.Create(context, account)

		assert.Nil(t, error)
		assert.NotNil(t, account)
	})

}
func TestAccounts_Fetch(t *testing.T) {

	t.Run("fetches account with valid data", func(*testing.T) {
		client := NewClient()
		context := context.Background()

		account := &Account{
			Data: &AccountData{
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Type:           "accounts",
				Attributes: &AccountAttributes{
					Country: "GB",
					Name:    []string{"John Doe"},
				},
			},
		}

		account, _ = client.Accounts.Create(context, account)
		fetchedAccount, error := client.Accounts.Fetch(context, account.Data.ID)

		assert.Nil(t, error)
		assert.NotNil(t, fetchedAccount)
	})
}
