//go:build integration

package form3_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/castanhojfc/form3-client-go/form3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	suite.databaseConnection.Session(&gorm.Session{AllowGlobalUpdate: true}).Table("public.Account").Delete("true")
}

type TestCase struct {
	description string
	request     string
	expected    string
}

func (suite *Form3AccountsTestSuite) Test_Create() {
	suite.T().Parallel()

	suite.T().Run("should create account when a valid account is provided", func(t *testing.T) {
		client, _ := form3.New()

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

				account, response, error := client.Accounts.Create(account)

				assert.Equal(t, account, expected)
				assert.NotNil(t, response)
				assert.Nil(t, error)
			})
		}
	})

	suite.T().Run("should not create account when there is a problem marshalling the acount", func(t *testing.T) {
		mockJsonMarshal := new(JsonMarshalMock)
		mockJsonMarshal.On("Marshal", mock.Anything).Return(nil, fmt.Errorf("marshalling issue"))

		client, _ := form3.New()
		client.Accounts.JsonMarshal = mockJsonMarshal.Marshal

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "d3f29952-ab3b-4dc3-bc1e-adbb6e1ff98e"
		account, response, error := client.Accounts.Create(account)

		assert.Equal(suite.T(), error, form3.OperationError{Message: "marshalling issue", Body: nil})
		assert.Nil(suite.T(), response)
		assert.Nil(suite.T(), account)
		mockJsonMarshal.AssertExpectations(t)
	})

	suite.T().Run("should not create account when there is a problem perfoming the request", func(t *testing.T) {
		client, _ := form3.New()
		client.BaseUrl = &url.URL{
			Scheme: "asdf",
			Host:   "asdf",
		}

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "0027c3aa-3aa4-4306-9efa-4b8472d875c1"
		account, response, error := client.Accounts.Create(account)

		assert.Contains(suite.T(), error.Error(), "unsupported protocol scheme")
		assert.Nil(suite.T(), response)
		assert.Nil(suite.T(), account)
	})

	suite.T().Run("should not create account when unmarshalling and there is a reading body problem", func(t *testing.T) {
		mockReadAll := new(ReadAllMock)
		mockReadAll.On("ReadAll", mock.Anything).Return(nil, fmt.Errorf("read issue"))

		client, _ := form3.New()
		client.Accounts.ReadAll = mockReadAll.ReadAll

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "8a3f59a4-7d55-400b-b561-1eb6b68ad8fa"
		account, response, error := client.Accounts.Create(account)

		assert.Equal(suite.T(), form3.OperationError{Message: "read issue", Body: nil}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
		mockReadAll.AssertExpectations(t)
	})

	suite.T().Run("should not create account when unmarshalling and there is a unmarshalling problem", func(t *testing.T) {
		mockJsonUnmarshal := new(JsonUnmmarshalMock)
		mockJsonUnmarshal.On("Unmarshal", mock.Anything, mock.Anything).Return(fmt.Errorf("unmarshal issue"))

		client, _ := form3.New()
		client.Accounts.JsonUnmarshal = mockJsonUnmarshal.Unmarshal

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "796a9db8-6159-46c8-8f78-9be07c93c24c"
		account, response, error := client.Accounts.Create(account)

		assert.Equal(suite.T(), form3.OperationError{Message: "unmarshal issue", Body: nil}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
		mockJsonUnmarshal.AssertExpectations(t)
	})

	suite.T().Run("should not create account when an account without required information is provided", func(*testing.T) {
		client, _ := form3.New()

		account := accountFromJson(suite.T(), "./fixtures/requests/account_missing_required_data.json")

		client.Accounts.Create(account)
		account.Data.ID = "c0582554-867d-42d3-a62e-1d64ae9f5b8e"
		account, response, error := client.Accounts.Create(account)

		assert.Equal(suite.T(), form3.OperationError{Message: "400 Bad Request", Body: []byte("{\"error_message\":\"validation failure list:\\nvalidation failure list:\\norganisation_id in body is required\"}")}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
	})

	suite.T().Run("should not create account when an account was previously created", func(*testing.T) {
		client, _ := form3.New()

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "ab7278a5-9c8e-4760-b69a-6f83b73e1b53"

		client.Accounts.Create(account)
		account, response, error := client.Accounts.Create(account)

		assert.Equal(suite.T(), form3.OperationError{Message: "409 Conflict", Body: []byte("{\"error_message\":\"Account cannot be created as it violates a duplicate constraint\"}")}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
	})
}

func (suite *Form3AccountsTestSuite) Test_Fetch() {
	suite.T().Parallel()

	suite.T().Run("should fetch an existing account", func(*testing.T) {
		client, _ := form3.New()

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "999a01ef-2695-48f0-b6b6-54c8a30faa3f"

		account, _, _ = client.Accounts.Create(account)
		fetchedAccount, response, error := client.Accounts.Fetch(account.Data.ID)

		assert.Nil(suite.T(), error)
		assert.Equal(suite.T(), account, fetchedAccount)
		assert.NotNil(suite.T(), response)
	})

	suite.T().Run("should not fetch account when there is a problem perfoming the request", func(t *testing.T) {
		client, _ := form3.New()
		client.BaseUrl = &url.URL{
			Scheme: "asdf",
			Host:   "asdf",
		}

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "57238e6f-fc28-4d63-8e31-d901882b104f"

		client.Accounts.Create(account)
		fetchedAccount, response, error := client.Accounts.Fetch(account.Data.ID)

		assert.Contains(suite.T(), error.Error(), "unsupported protocol scheme")
		assert.Nil(suite.T(), response)
		assert.Nil(suite.T(), fetchedAccount)
	})

	suite.T().Run("should not fetch account when the account is does not exist", func(t *testing.T) {
		client, _ := form3.New()

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "f65b0db1-50b9-4ef3-81b4-1a9442d75d0c"
		account, response, error := client.Accounts.Fetch(account.Data.ID)

		assert.Equal(suite.T(), form3.OperationError{Message: "404 Not Found", Body: []byte("{\"error_message\":\"record f65b0db1-50b9-4ef3-81b4-1a9442d75d0c does not exist\"}")}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
	})

	suite.T().Run("should not fetch account when there is a problem perfoming the request", func(t *testing.T) {
		client, _ := form3.New()
		client.BaseUrl = &url.URL{
			Scheme: "asdf",
			Host:   "asdf",
		}

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "26eeb841-edd5-4d9e-947f-db60f91a7085"
		account, response, error := client.Accounts.Fetch(account.Data.ID)

		assert.Contains(suite.T(), error.Error(), "unsupported protocol scheme")
		assert.Nil(suite.T(), response)
		assert.Nil(suite.T(), account)
	})

	suite.T().Run("should not fetch account when unmarshalling and there is a reading body problem", func(t *testing.T) {
		mockReadAll := new(ReadAllMock)
		mockReadAll.On("ReadAll", mock.Anything).Return(nil, fmt.Errorf("read issue"))

		client, _ := form3.New()
		client.Accounts.ReadAll = mockReadAll.ReadAll

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "bf81ac45-3b70-4ec9-946e-ec9d4b651b0d"

		client.Accounts.Create(account)
		account, response, error := client.Accounts.Fetch(account.Data.ID)

		assert.Equal(suite.T(), form3.OperationError{Message: "read issue", Body: nil}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
		mockReadAll.AssertExpectations(t)
	})

	suite.T().Run("should not fetch account when unmarshalling and there is a unmarshalling problem", func(t *testing.T) {
		mockJsonUnmarshal := new(JsonUnmmarshalMock)
		mockJsonUnmarshal.On("Unmarshal", mock.Anything, mock.Anything).Return(fmt.Errorf("unmarshal issue"))

		client, _ := form3.New()
		client.Accounts.JsonUnmarshal = mockJsonUnmarshal.Unmarshal

		var account = accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "ae8332af-2256-49de-adb7-e1c596430c8e"

		client.Accounts.Create(account)
		account, response, error := client.Accounts.Fetch(account.Data.ID)

		assert.Equal(suite.T(), form3.OperationError{Message: "unmarshal issue", Body: nil}, error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), account)
		mockJsonUnmarshal.AssertExpectations(t)
	})
}

func (suite *Form3AccountsTestSuite) Test_Delete() {
	suite.T().Parallel()

	suite.T().Run("should delete account when the account exists", func(*testing.T) {
		client, _ := form3.New()

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "cf8a82a8-376f-4572-9cc4-e73578cf99e7"

		client.Accounts.Create(account)
		response, error := client.Accounts.Delete(account.Data.ID, 0)
		fetchedAccount, _, _ := client.Accounts.Fetch(account.Data.ID)

		assert.Nil(suite.T(), error)
		assert.NotNil(suite.T(), response)
		assert.Nil(suite.T(), fetchedAccount)
	})

	suite.T().Run("should not delete account when there is a problem performing the request", func(*testing.T) {
		client, _ := form3.New()

		client.BaseUrl = &url.URL{
			Scheme: "asdf",
			Host:   "asdf",
		}

		account := accountFromJson(suite.T(), "./fixtures/requests/uk_account_with_confirmation_of_payee.json")
		account.Data.ID = "b0a7d0e2-ca99-42de-8655-1e4ff0794cb2"

		client.Accounts.Create(account)
		response, error := client.Accounts.Delete(account.Data.ID, 0)

		assert.Contains(suite.T(), error.Error(), "unsupported protocol scheme")
		assert.Nil(suite.T(), response)
	})

	suite.T().Run("should not delete account when there the account is not existent", func(t *testing.T) {
		client, _ := form3.New()

		response, error := client.Accounts.Delete("5faad046-ca12-475b-be4e-425c9668d3ab", 0)

		assert.Equal(suite.T(), form3.OperationError{Message: "404 Not Found", Body: []byte{}}, error)
		assert.NotNil(suite.T(), response)
	})

	suite.T().Run("should not delete account when there is a reading body problem", func(t *testing.T) {
		mockReadAll := new(ReadAllMock)
		mockReadAll.On("ReadAll", mock.Anything).Return(nil, fmt.Errorf("read issue"))

		client, _ := form3.New()
		client.Accounts.ReadAll = mockReadAll.ReadAll

		response, error := client.Accounts.Delete("5fafd046-sd42-475b-be4e-425c5468d3ab", 0)

		assert.Equal(suite.T(), form3.OperationError{Message: "read issue", Body: nil}, error)
		assert.NotNil(suite.T(), response)
		mockReadAll.AssertExpectations(t)
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

type JsonMarshalMock struct {
	mock.Mock
}

func (m *JsonMarshalMock) Marshal(v any) ([]byte, error) {
	args := m.Called(v)

	return []byte{}, args.Error(1)
}

type JsonUnmmarshalMock struct {
	mock.Mock
}

func (m *JsonUnmmarshalMock) Unmarshal(data []byte, v any) error {
	args := m.Called(data, v)

	return args.Error(0)
}

type ReadAllMock struct {
	mock.Mock
}

func (m *ReadAllMock) ReadAll(r io.Reader) ([]byte, error) {
	args := m.Called(r)

	return []byte{}, args.Error(1)
}
