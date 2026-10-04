.PHONY: fmt checkfmt vet lint tidy check build run

# Format all Go files in place
fmt:
	gofmt -w .

# Fail if any file is not gofmt-clean (no changes made)
checkfmt:
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then echo "not gofmt-clean:"; echo "$$files"; exit 1; fi

# Static analysis (stdlib)
vet:
	go vet ./...

# Full linter (golangci-lint)
lint:
	golangci-lint run ./...

# Sync go.mod / go.sum with actual imports
tidy:
	go mod tidy

# Everything: format + all checks
check: fmt vet lint tidy

# Build binary (same as your_program.sh)
build:
	go build -o /tmp/codecrafters-build-claude-code-go app/*.go

# Run with a prompt: make run PROMPT="hello"
run:
	go run app/*.go -p "$(PROMPT)"
