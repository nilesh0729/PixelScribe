Container:
	docker run --name Steno -p 5434:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:17-alpine

CreateDB:
	docker exec -it Steno createdb --username=root --owner=root Pros

DropDB:
	docker exec -it Steno dropdb -U root Pros

DB_URL=postgresql://root:secret@localhost:5434/Pros?sslmode=disable

MigrateUp:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

MigrateDown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

Sqlc:
	sqlc generate

Test:
	go test -v -cover ./...

.PHONY:	Container	CreateDB	DropDB	MigrateUp	MigrateDown	Sqlc	Test