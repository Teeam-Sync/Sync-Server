dev:
	export APP_ENV=dev && go run ./cmd/main.go

prod:
	export APP_ENV=prod && go run ./cmd/main.go