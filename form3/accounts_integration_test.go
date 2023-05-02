//go:build integration

package form3_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/castanhojfc/form3-client-go/form3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Form3AccountsTestSuite struct {
	suite.Suite
	databaseConnection *gorm.DB
}

func TestForm3AccountsTestSuite(t *testing.T) {
	suite.Run(t, new(Form3AccountsTestSuite))
}

func (suite *Form3AccountsTestSuite) SetupTest() {
	host := os.Getenv("TEST_DATABASE_HOST")
	user := os.Getenv("TEST_DATABASE_USERNAME")
	password := os.Getenv("TEST_DATABASE_PASSWORD")
	dbname := os.Getenv("TEST_DATABASE_NAME")
	port := os.Getenv("TEST_DATABASE_PORT")
	sslmode := os.Getenv("TEST_DATABASE_SSL_MODE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	database, error := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if error != nil {
		suite.T().Fatalf("Failed to connect to the test database")
	}

	suite.databaseConnection = database

	suite.databaseConnection.Exec("DELETE FROM \"public\".\"Account\"")
}

type TestCase struct {
	description string
	request     string
	expected    string
}

func (suite *Form3AccountsTestSuite) Test_Create() {
	suite.T().Parallel()

	suite.T().Run("should create account when a valid account is provided", func(t *testing.T) {
		client, _ := form3.NewClient()

		tests := []TestCase{
			{
				description: "UK account with confirmation of payee",
				request:     "./fixtures/requests/uk_account_with_confirmation_of_payee.json",
				expected:    "./fixtures/responses/uk_account_with_confirmation_of_payee.json",
			},
			{
				description: "UK account without confirmation of payee",
				request:     "./fixtures/requests/uk_account_without_confirmation_of_payee.json",
				expected:    "./fixtures/responses/uk_account_without_confirmation_of_payee.json",
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				account := accountFromJson(t, test.request)
				expected := accountFromJson(t, test.expected)

				account, _ = client.Accounts.Create(account)

				assert.Equal(t, expected, account)
			})
		}
	})

	suite.T().Run("should not create account when an account without required information is provided", func(*testing.T) {
		client, _ := form3.NewClient()

		account := accountFromJson(suite.T(), "./fixtures/requests/account_missing_required_data.json")

		client.Accounts.Create(account)
		account, error := client.Accounts.Create(account)

		assert.True(suite.T(), strings.Contains(fmt.Sprint(error), "organisation_id in body is required"))
		assert.Nil(suite.T(), account)
	})

	suite.T().Run("should not create account when an account was previously created", func(*testing.T) {
		client, _ := form3.NewClient()

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")

		client.Accounts.Create(account)
		account, error := client.Accounts.Create(account)

		assert.True(suite.T(), strings.Contains(fmt.Sprint(error), "Account cannot be created as it violates a duplicate constraint"))
		assert.Nil(suite.T(), account)
	})
}

func (suite *Form3AccountsTestSuite) Test_Fetch() {
	suite.T().Parallel()

	suite.T().Run("should fetch an existing account", func(*testing.T) {
		client, _ := form3.NewClient()

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = uuid.New().String()

		client.Accounts.Create(account)
		fetchedAccount, error := client.Accounts.Fetch(account.Data.ID)

		assert.Nil(suite.T(), error)
		assert.NotNil(suite.T(), fetchedAccount)
	})
}

func (suite *Form3AccountsTestSuite) Test_Delete() {
	suite.T().Parallel()

	suite.T().Run("should delete account when the account exists", func(*testing.T) {
		client, _ := form3.NewClient()

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = uuid.New().String()

		client.Accounts.Create(account)
		error := client.Accounts.Delete(account.Data.ID, 0)
		fetchedAccount, _ := client.Accounts.Fetch(account.Data.ID)

		assert.Nil(suite.T(), error)
		assert.Nil(suite.T(), fetchedAccount)
	})
}

func accountFromJson(t *testing.T, fileName string) *form3.Account {
	file, error := os.Open(fileName)

	if error != nil {
		t.Fatalf("failed to open fixture file: %v", error)
	}

	defer file.Close()

	fileBytes, error := ioutil.ReadAll(file)

	if error != nil {
		t.Fatalf("failed to read fixture bytes: %v", error)
	}

	account := &form3.Account{}
	json.Unmarshal(fileBytes, &account)

	return account
}
