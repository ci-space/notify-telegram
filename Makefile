test: ## Run tests
	go test ./...
	cd ./pkg/md2html && go test ./...
	cd ./pkg/tgapi && go test ./...

lint: ## Run lint
	golangci-lint run --fix

build/docker:
	docker build -t ci-space/notify-telegram:local .
