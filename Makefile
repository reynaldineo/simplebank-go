postgres:
	docker run --name postgres-alpine -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:alpine

createdb: 
	docker exec -it postgres-alpine createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-alpine dropdb simple_bank

migrateup:
	go run db/main.go -action up

migrateup1:
	go run db/main.go -action up1

migratedown:
	go run db/main.go -action down

migratedown1:
	go run db/main.go -action down1

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc