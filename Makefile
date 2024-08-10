.PHONY: install-tools check test

# Go tool paths
GOLINT = $(shell go env GOPATH)/bin/golint
INEFFASSIGN = $(shell go env GOPATH)/bin/ineffassign
MISSPELL = $(shell go env GOPATH)/bin/misspell
GOCYCLO = $(shell go env GOPATH)/bin/gocyclo

install-tools:
	@echo "Installing tools..."
	go install golang.org/x/lint/golint@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

check: install-tools
	@echo "Running checks..."
	go fmt ./...
	go vet ./...
	$(GOLINT) ./...
	$(MISSPELL) -w .
	$(GOCYCLO) -over 10 .
	$(INEFFASSIGN) .

test:
	@echo "Running tests..."
	go test -coverprofile=coverage.out -coverpkg=$$(go list ./... | grep -v /test$$ | grep -v main | grep -v '_repository.go$$' | tr '\n' ',') ./...

coverage-report: test
	@echo "Generating coverage report..."
	go tool cover -func=coverage.out | grep total | awk '{print substr($$NF, 1, length($$NF)-1)}' > coverage.txt
	@COVERAGE=$$(cat coverage.txt); \
	echo "Coverage: $$COVERAGE"; \
	if (( $$(echo "$$COVERAGE < 30" | bc -l) )); then \
		echo "Coverage is below 30%. Stopping the process."; \
		exit 1; \
	fi
