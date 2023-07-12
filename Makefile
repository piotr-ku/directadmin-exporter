test:
	@go test -cover ./...
	@rm -f coverage.out

coverage:
	@go test ./... -coverprofile=./coverage.out 2> /dev/null
	@go tool cover -html=coverage.out
	@go test ./... -coverprofile=./coverage.out -covermode=atomic
	@${GOPATH}/bin/go-test-coverage --config=./.testcoverage.yml
	@rm -f coverage.out