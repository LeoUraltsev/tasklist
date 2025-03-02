run_app:
	go run ./cmd/tasklist/main.go

goose_create_migrations_user:
	goose -dir migrations create user_table sql

goose_up:
	goose sqlite3 -dir migrations ./storage/tasklist.db up

goose_status:
	goose sqlite3 -dir migrations ./storage/tasklist.db status

goose_down:
	goose sqlite3 -dir migrations ./storage/tasklist.db down

lint_run:
	 golangci-lint run