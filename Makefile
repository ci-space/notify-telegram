test: ## Run tests
	go test ./...

lint: ## Run lint
	golangci-lint run --fix
