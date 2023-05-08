# Hi there 👋 thank you for taking a look at my code!

I'm Francisco Castanho and I am new to golang.

# How to run all tests 🧪

```
docker-compose up
```

After they're all executed, package documentation is served using godocs using a web server. 📖

[Checkout this link once the documentation is served](http://localhost:6060/pkg/github.com/castanhojfc/form3-client-go/form3/)

# Run bare metal 🤘

Once the account API is running via docker compose, one can attach to the container and execute the commands provided by the Makefile.

To run all commands without attaching, just run these commands to make all environment variables available and the api reachable from the host machine:

```
echo "127.0.0.1 accountapi" | sudo tee -a /etc/hosts
source local.env
```

Here are the most important ones:
- Run all checks `make checks`.
- Generate a coverage report just run `make generate_report` and open the generated report using a browser. (The coverage is 100% ✅ which is nice, although that does not mean that they are perfect 😆)
- Generate all documentation and serve it via a web server `make generate_docs`

# How to use 🧑‍💻

The default base URL is `http://accountapi:8080`.

```
import "github.com/castanhojfc/form3-client-go/form3"

// Create a API client
client, error := form3.New()

// Set options, there are defaults already setup
client.BaseUrl = url.ParseRequestURI("http://asdf:8080")
client.HttpClient = &http.Client{}
client.DebugEnabled = true
client.HttpRetryAttempts = 4
client.HttpTimeUntilNextAttempt = 3 * time.Second
client.HttpTimeout = 10 * time.Second

// Build an account object
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

// Create an acount
account, response, error = client.Accounts.Create(account)

// Fetch an account, takes the account id as an argument
account, response, error := client.Accounts.Fetch("5e759a85-e632-4b5d-8232-494552d11212")

// Delete an account, takes the account id and version as arguments
response, error := client.Accounts.Delete("5e759a85-e632-4b5d-8232-494552d11212", 0)
```

In all operations, a HTTP request is returned if successfully performed.

This is so that the caller can inspect exactly what happened, even if later on another error occurs.

The client should be able to handle retries when there's a chance of making a successful request in the future. Additionally it should be able to handle client timeouts and make it self identifiable to the server.

More details in the docs! 📖

## Future work/Limitations 👷
 - More unit tests could have been written! I gave priority to integration tests.
 - Some tests could probably be table driven. I prioritized coverage and test quality.
 - Rate limit headers can be used to make more intelligent retry attempts.
 - There's no existence of tests checking the fields `created_on` and `modified_on` or even any other response coming from the server that shows a timestamp. This is because I was not able to freeze these dates.
 - I've used gock to mock http requests. Unfortunately it is not possible to run these tests in parallel, there must be a way to achieve this, but I was not able to.
 - To mock function calls from the standard library I´ve used dependency injection. Some parameters from the client and the account service exist and can be injected just for testing purposes.

## Bonus 🥳
  - Also added a continuous integration pipeline using GitHub Actions. It is simple, but its something.

# Form3 Take Home Exercise Original Description

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
