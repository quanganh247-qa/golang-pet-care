DB_URL=postgresql://admin:MW0I7JWVOfRYTKstGCmJS7IO7IKELVBH@dpg-cqvkhntds78s739lgfag-a.virginia-postgres.render.com/blog_8g7e
name= init_db
mup:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose up
mdown:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose down
new_migration:
	migrate create -ext sql -dir app/db/migration -seq $(name)
mforce:
	migrate -path app/db/migration -database "$(DB_URL)" -verbose force 1
migrateup-github:
	migrate -path app/db/migration -database "postgresql://admin:MW0I7JWVOfRYTKstGCmJS7IO7IKELVBH@dpg-cqvkhntds78s739lgfag-a.virginia-postgres.render.com/blog_8g7e" -verbose up
	 
sqlc:
	docker run --rm -v ".://src" -w //src sqlc/sqlc:1.20.0 generate 

test:
	go test -v -cover ./...
	
server:
	gin -p 8081 -i run main.go

mock:
	mockgen -package mockdb -destination app/db/mock/store.go github.com/hpt/go-client/app/db/sqlc Store

proto:
	rm -rf app/pb/*.go
	protoc --proto_path=app/proto  --go_out=app/pb --go_opt=paths=source_relative \
	--go-grpc_out=app/pb --go-grpc_opt=paths=source_relative \
	app/proto/*.proto 

evans:
	evans --host localhost --port 9090 -r repl

postgres:
	docker run -d  --name postgres  -p 5432:5432 -e POSTGRES_PASSWORD=12345678  -e PGDATA=/var/lib/postgresql/data/pgdata  -v postgres_volume:/var/lib/postgresql/data  postgres:15-alpine

redis:
	docker run -d --name redis -p 6379:6379 redis:7-alpine

.PHONY: mup mdown new_migration mforce sqlc test server mock proto evans postgres migrateup-github redis mup_test mdown_test mforce_test