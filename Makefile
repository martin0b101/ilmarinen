DB_URL=postgres://admin:password@0.0.0.0:5432/ilmarinen?sslmode=disable

.PHONY: test database-up database-down migrate-up migrate-down


default: 
	$(MAKE) fmt
	$(MAKE) tidy
	$(MAKE) lint
	$(MAKE) test

test:
	go test ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

lint:
	go vet ./...

database-up:
	docker componse up -d

database-down:
	docker compone down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1
