format:
	gofmt -s -w -l .

lint:
	golangci-lint run

test:
	go test -race -shuffle=on -parallel 16 -v ./... -tags=unit,integration

integration_test:
	go test -race -shuffle=on -parallel 16 -v ./... -tags=integration

coverage:
	go test ./... --cover

generate_report: coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

check: format lint test