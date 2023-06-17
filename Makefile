.PHONY: help
help:
	@echo List of commands:
	@echo   test                    - run tests
	@echo   docker-up               - docker compose up
	@echo   docker-down             - docker compose down
	@echo Usage:
	@echo          make `cmd_name`

.PHONY: test
test:
	go test -cover ./...

.PHONY: docker-up
docker-up:
	docker compose up

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: build-cli
build-cli:
	go build -C ./cmd/cli longsub