.PHONY: run-worker-service
run-worker-service:
	go run ./cmd/main.go worker

.PHONY: run-mail-app-service
run-mail-app-service:
	go run ./cmd/main.go app

.PHONY: run-migrate
run-migrate:
	go run ./cmd/main.go migrate