include .env
export

run:
	migrate -path ./migrations -database "mysql://root:@tcp(localhost:3306)/go_clean?parseTime=true" -verbose up
	go run ./cmd/api/main.go

build:
	go build ./cmd/api/main.go

test:
	go test -cover ./...


down:
	migrate -path ./migrations -database "mysql://root:@tcp(localhost:3306)/go_clean?parseTime=true" -verbose down