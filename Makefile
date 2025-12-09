.PHONY: docs
docs:
	cd backend && swag init -g cmd/main.go -o docs --generatedTime=false --parseDependency --parseInternal

.PHONY: run_containers
run_containers:
	docker-compose up

.PHONY: stop_containers
stop_containers:
	docker-compose down -v

.PHONY: run_server
run_server:
	cd backend && go run ./cmd/main.go

.PHONY: fmt
fmt:
	@echo "🧹 Formatting Go code..."
	@gofmt -l -w `find . -type f -name '*.go' -not -path "./vendor/*"`
	@golines --max-len=120 --base-formatter=gofmt --shorten-comments --ignore-generated  --ignored-dirs=vendor -w .
	@echo "✅ Code formatted successfully"

.PHONY: run_frontend
run_frontend:
	cd frontend && npm run dev -- --host