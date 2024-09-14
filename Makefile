test:
	@go test -v ./tests/...

coverage:
	@go test -coverprofile=coverage.out ./parsex/... ./tests/...
	@go tool cover -html=coverage.out -o ./tests/coverage.html
	@rm coverage.out

push: coverage
	git add ./tests/coverage.html
	git commit -m "tests(make): Update test coverage"
	git push -u origin main
