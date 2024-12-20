migrate_up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate_up_next:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate_down_last:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

db_up:
	- docker pull postgres:12-alpine
	- docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
	- sleep 2
	- docker exec -it postgres12 createdb --username=root --owner=root simple_bank

db_down:
	- docker stop postgres12
	- docker remove postgres12
	- docker rmi postgres:12-alpine

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination=db/mock/store.go -package=mockdb github.com/T-BO0/bank/db/sqlc Store  

.PHONY: db_up db_down migrate_up migrate_down sqlc test server mock migrate_up_next migrate_down_last