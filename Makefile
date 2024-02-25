local:
	go run cmd/app/main.go --config=./config/local.yaml

dev:
	go get -u .../
	go run cmd/app/main.go --config=./config/dev.yaml

prod:
	go get -u .../
	go run cmd/app/main.go --config=./config/prod.yaml

script-migrations:
	go run ./cmd/migration --migrations-path=./migrations