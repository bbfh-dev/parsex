test:
	@go test -v ./...

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o ./tests/coverage.html
	@rm coverage.out
