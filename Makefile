format:
	gofmt -s -w -l .

lint:
	golangci-lint run

test:
	go test -v ./...

coverage:
	go test ./... --cover

report: coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

check: format lint test