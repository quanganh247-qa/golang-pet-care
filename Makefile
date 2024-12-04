DB_URL=postgresql://postgres:12345678@localhost:5433/postgres?sslmode=disable
mup:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose up
mdown:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose down
new_migration:
	migrate create -ext sql -dir app/db/migration -seq change_table
mforce:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose force 1
migrateup-github:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose up
	 
sqlc:
	docker run --rm -v ".://src" -w //src sqlc/sqlc:1.20.0 generate 

test:
	go test -v -cover ./...
	
server:
	gin -p 8089 -i run main.go


postgres:
	docker run -d  --name postgres  -p 5433:5432 -e POSTGRES_PASSWORD=12345678  -e PGDATA=/var/lib/postgresql/data/pgdata  -v postgres_volume:/var/lib/postgresql/data  postgres:15-alpine

redis:
	docker run -d --name redis -p 6379:6379 redis:7-alpine

rabbitmq:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

postgres-ec:
	docker run -d  --name postgres-ec  -p 5432:5432 -e POSTGRES_PASSWORD=12345678  -e PGDATA=/var/lib/postgresql/data/pgdata  -v postgres_volume:/var/lib/postgresql/data  postgres:15-alpine


.PHONY: mup mdown  mforce sqlc server   postgres  redis  rabbitmq