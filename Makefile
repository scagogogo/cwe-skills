.PHONY: build test test-ci coverage vet fmt lint clean examples

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

# Coverage parameters
COVERAGE_DIR=coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Build
build:
	$(GOBUILD) ./...

# Test
test:
	$(GOTEST) -v -count=1 ./...

# Test with race detector
test-race:
	$(GOTEST) -v -race -count=1 ./...

# Test for CI (with coverage)
test-ci:
	$(GOTEST) -v -race -count=1 -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Coverage report
coverage:
	$(GOTEST) -v -count=1 -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

# Coverage HTML report
coverage-html:
	$(GOTEST) -v -count=1 -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Vet
vet:
	$(GOVET) ./...

# Format
fmt:
	$(GOFMT) ./...

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

# Module tidy
tidy:
	$(GOMOD) tidy

# Clean
clean:
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -rf $(COVERAGE_DIR)

# Examples
examples:
	@echo "Running examples..."
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Running $$dir"; \
			cd "$$dir" && $(GOBUILD) -o /tmp/example_binary ./ && /tmp/example_binary && cd -; \
		fi; \
	done

# All checks
check: vet fmt test coverage