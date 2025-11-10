Container:
	docker run --name Steno -p 5434:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:18.0-alpine3.22

CreateDB:
	docker exec -it Steno createdb --username=root --owner=root Pros

DropDB:
	docker exec -it Steno dropdb -U root Pros

MigrateUp:
	migrate -path db/migration -database "postgres://root:secret@localhost:5434/Pros?sslmode=disable" -verbose up

MigrateDown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5434/Pros?sslmode=disable" -verbose down

.PHONY:	Container	CreateDB	DropDB	MigrateUp	MigrateDown