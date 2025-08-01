.PHONY: sqlc migrateup migratedown test server build clean

# Generate SQL code
sqlc:
	docker-compose run --rm api sh -c "sqlc generate"

# Run database migrations up
migrateup:
	docker-compose run --rm api sh -c "migrate -path db/migration -database 'postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable' -verbose up"

# Run database migrations down
migratedown:
	docker-compose run --rm api sh -c "migrate -path db/migration -database 'postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable' -verbose down"

# Run tests with coverage
test:
	make migrateup
	docker-compose run --rm api sh -c "go test -v -cover ./..."
	make migratedown


run:
	docker-compose up --build

destroy:
	docker-compose down
	make remove-data


# Remove postgres data volume
remove-data:
	docker volume rm postgres_data 2>/dev/null || true
