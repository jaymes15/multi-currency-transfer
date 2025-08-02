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
	docker-compose up -d
	sleep 3
	make migrateup
	@echo "🧪 Running all tests..."
	docker-compose run --rm api sh -c "go test -cover ./..."
	@echo "✅ All tests completed!"

test-verbose:
	make migrateup
	@echo "🧪 Running all tests with verbose output..."
	docker-compose run --rm api sh -c "go test -v -cover ./..."
	make migratedown


test-db:
	@echo "🧪 Running database tests only..."
	docker-compose run --rm api sh -c "go test -cover ./db/sqlc/"

test-clean:
	@echo "🧪 Running tests with clean output..."
	./scripts/test.sh

test-single:
	make migrateup
	docker-compose run --rm api sh -c "go test -v -cover $(TEST) || exit 1"
	make migratedown
run:
	docker-compose up --build

destroy:
	docker-compose down
	make remove-data


# Remove postgres data volume
remove-data:
	docker volume rm postgres_data 2>/dev/null || true

mockgen:
	docker-compose run --rm api sh -c "cd /app && mockgen -package mockdb -destination db/mock/store.go lemfi/simplebank/db/sqlc Store"