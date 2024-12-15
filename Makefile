DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version slate-up slate-down slate-nuke

migrate-up:
	migrate -database "${DB_URL}" -path db/migrations up

migrate-down:
	migrate -database "${DB_URL}" -path db/migrations down

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-force:
	migrate -database "${DB_URL}" -path db/migrations force $(version)

migrate-version:
	migrate -database "${DB_URL}" -path db/migrations version

slate-up:
	docker compose up -d

slate-down:
	docker compose down

slate-nuke:
	docker compose down -v