export

# The output the help for each task
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run-local: ### Run backend locally
	go run cmd/app/main.go
.PHONY: run

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable/bidding' up
.PHONY: migrate-up

migrate-create: ### migrate create new migration (up and down)
	migrate create -ext sql -dir migrations -seq bids
.PHONY: migrate-create

