# Hi there 👋 thank you for taking a look at my code!

I'm Francisco Castanho and I am new to golang

# How to run all tests

```
docker-compose up
```

Documentation using godocs is generated and available [here](http://localhost:6060/pkg/github.com/castanhojfc/form3-client-go/form3/)

# How to use

The base URL must be accessible. It should be configured in the environment variable `FORM3_ACCOUNT_API_URL` or as an option.

## Import the package
```
import "github.com/castanhojfc/form3-client-go/form3"
```

## Create a client
```
client, error := form3.New()
```

## (Optional) Configure the client when creating a client
```
// It is not needed to provide all options
url, _ := url.ParseRequestURI("http://asdf:8080")
httpClient := &http.Client{}

client, error := form3.New(
  form3.WithBaseUrl(url), // If both the option is provided and the environment variable, the option is honored
  form3.WithHttpClient(&http.Client{})
)
```

## Create an account
```
account := &form3.Account{
  Data: &form3.AccountData{
    ID:             "47cf8708-3c26-4baa-b3d3-6365996e27c3",
    OrganisationID: "afe81b33-210b-42a5-8d80-40e5adde721e",
    Type:           "accounts",
    Attributes: &form3.AccountAttributes{
      Country:                 "GB",
      BaseCurrency:            "GBP",
      BankID:                  "400302",
      BankIDCode:              "GBDSC",
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

account, error = client.Accounts.Create(account)
```

## Fetch an account
```
// Fetch takes the account id as an argument
account, error := client.Accounts.Fetch("5e759a85-e632-4b5d-8232-494552d11212")
```

## Delete an account
```
// Delete takes the account id and version as arguements
client.Accounts.Delete("5e759a85-e632-4b5d-8232-494552d11212", 0)
```

## Observations
  - There is a makefile with a few useful commands available. Check it out! :partying_face:
  - Documentation is available through godocs.
  - Also added a continuous integration pipeline using GitHub Actions as a bonus.

## Future work/Limitations
 - More unit could have been written! I gave priority to integration tests.
 - Could not make DumpRequest return an error to cover more logic using an integration test.

# Form3 Take Home Exercise

Engineers at Form3 build highly available distributed systems in a microservices environment. Our take home test is designed to evaluate real world activities that are involved with this role. We recognise that this may not be as mentally challenging and may take longer to implement than some algorithmic tests that are often seen in interview exercises. Our approach however helps ensure that you will be working with a team of engineers with the necessary practical skills for the role (as well as a diverse range of technical wizardry). 

## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](https://www.api-docs.form3.tech/api/tutorials/getting-started/create-an-account) for information on how to interact with the API. Please note that the fake account API does not require any authorisation or authentication.

A mapping of account attributes can be found in [models.go](./models.go). Can be used as a starting point, usage of the file is not required.

If you encounter any problems running the fake account API we would encourage you to do some debugging first,
before reaching out for help.

## Submission Guidance

### Shoulds

The finished solution **should:**
- Be written in Go.
- Use the `docker-compose.yaml` of this repository.
- Be a client library suitable for use in another software project.
- Implement the `Create`, `Fetch`, and `Delete` operations on the `accounts` resource.
- Be well tested to the level you would expect in a commercial environment. Note that tests are expected to run against the provided fake account API.
- Run the tests when `docker-compose up` is run - our reviewers will run `docker-compose up` and expect to see the test results in the output.
- Be simple and concise.

### Should Nots

The finished solution **should not:**
- Use a code generator to write the client library.
- Use (copy or otherwise) code from any third party without attribution to complete the exercise, as this will result in the test being rejected.
    - **We will fail tests that plagiarise others' work. This includes (but is not limited to) other past submissions or open-source libraries.**
- Use a library for your client (e.g: go-resty). Anything from the standard library (such as `net/http`) is allowed. Libraries to support testing or types like UUID are also fine.
- Implement client-side validation.
- Implement an authentication scheme.
- Implement support for the fields `data.attributes.private_identification`, `data.attributes.organisation_identification`
  and `data.relationships` or any other fields that are not included in the provided `models.go`, as they are omitted from the provided fake account API implementation.
- Have advanced features, however discussion of anything extra you'd expect a production client to contain would be useful in the documentation.
- Be a command line client or other type of program - the requirement is to write a client library.
- Implement the `List` operation.
> We give no credit for including any of the above in a submitted test, so please only focus on the "Shoulds" above.

## How to submit your exercise

- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, by copying all files you deem necessary for your submission
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) [@form3tech-interviewer-1](https://github.com/form3tech-interviewer-1) to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team

## License

Copyright 2019-2023 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
