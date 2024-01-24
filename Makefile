migrate-up:
	goose -env local -path internal/storage/migrations/postgres up

migrate-down:
	goose -env local -path internal/storage/migrations/postgres down

.PHONY: gen
gen:
	go generate ./...
