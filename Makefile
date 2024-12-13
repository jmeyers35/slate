DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: migrate-up migrate-down

migrate-up:
	migrate -database "${DB_URL}" -path db/migrations up

migrate-down:
	migrate -database "${DB_URL}" -path db/migrations down