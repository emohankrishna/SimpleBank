postgres:
	docker run --name simplebank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it simplebank createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it simplebank dropdb simple_bank


migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate
server:
	go run main.go

test:
	go test ./... -v  -cover

mock:
	mockgen -destination db/mock/store.go --build_flags=--mod=mod -package mockdb github.com/emohankrishna/Simplebank/db/sqlc Store

.PHONY:
	postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 test server mock