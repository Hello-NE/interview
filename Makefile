postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres:18

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb --username=root simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

test:
	go test -v -cover -count=1 ./...

server:
	go run main.go

.PHONY: server createdb dropdb postgres migrateup migratedown test