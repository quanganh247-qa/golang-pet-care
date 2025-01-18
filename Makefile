DB_URL=postgresql://postgres:12345678@localhost:5433/postgres?sslmode=disable
SUPABASE_URL=postgresql://postgres.prmzavhkqqthcwdnrkgt:postgres@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres
mup:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose up
	
mdown:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir app/db/migration -seq change_table
<<<<<<< HEAD

=======
mforce:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose force 1
>>>>>>> dff4498 (calendar api)
migrateup-github:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose up

migrateup-supabase:
	migrate -path app/db/migration -database "$(SUPABASE_URL)" -verbose up
	
sqlc:
	docker run --rm -v ".://src" -w //src sqlc/sqlc:1.20.0 generate 

test:
	go test -v -cover ./...
	
server:
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	gin -p 8089 -i run /Users/dhquanganh/Documents/TLCN/golang-pet-care/main.go
=======
	gin -p 8081 -i run main.go
=======
	gin -p 8089 -i run main.go
>>>>>>> c3c833d (login api)
=======
	gin -p 8089 -i run /Users/dhquanganh/Documents/TLCN/golang-pet-care/main.go
>>>>>>> 685da65 (latest update)

>>>>>>> 9d28896 (image pet)

postgres:
	docker run -d  --name postgres  -p 5433:5432 -e POSTGRES_PASSWORD=12345678  -e PGDATA=/var/lib/postgresql/data/pgdata  -v postgres_volume:/var/lib/postgresql/data  postgres:15-alpine

redis:
	docker run -d --name redis -p 6379:6379 redis:7-alpine

elasticsearch:
	docker run --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" -e "xpack.security.enabled=false" -d elasticsearch:8.12.0

<<<<<<< HEAD
minio:
	docker run -d --name minio -p 9000:9000 -e "MINIO_ACCESS_KEY=1View" -e "MINIO_SECRET_KEY=12345678" -v minio_data:/data minio/minio:latest server /data

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

migrate-up:
	migrate -path app/db/migration -database "$(SUPABASE_URL)" up

migrate-down:
	migrate -path app/db/migration -database "$(SUPABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir app/db/migration -seq $(name)

.PHONY: mup mdown mforce sqlc server postgres redis supertokens elasticsearch minio docker-build docker-up docker-down docker-logs migrate-up migrate-down migrate-create
=======
postgres-ec:
	docker run -d  --name postgres-ec  -p 5432:5432 -e POSTGRES_PASSWORD=12345678  -e PGDATA=/var/lib/postgresql/data/pgdata  -v postgres_volume:/var/lib/postgresql/data  postgres:15-alpine


.PHONY: mup mdown  mforce sqlc server   postgres  redis  rabbitmq
>>>>>>> 21608b5 (cart and order api)
