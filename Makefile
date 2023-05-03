tidy:
	go mod tidy

format:
	gofmt -s -w -l .

lint:
	golangci-lint run

test:
	go test -race -shuffle=on -parallel 16 -v ./... -tags=unit,integration

integration_test:
	go test -race -shuffle=on -parallel 16 -v ./... -tags=integration

coverage:
	go test ./... --cover -tags=unit,integration

generate_report: coverage
	go test ./... -coverprofile=coverage.out -tags=unit,integration
	go tool cover -html=coverage.out

check: tidy format lint test

generate_docs:
	godoc -http :6060
