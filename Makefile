
# Generate SQL code
.PHONY: sqlc
sqlc:
	docker-compose run --rm api sh -c "sqlc generate"

# Run database migrations up
.PHONY: migrateup
migrateup:
	docker-compose run --rm api sh -c "migrate -path db/migration -database 'postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable' -verbose up"

# Run database migrations down
.PHONY: migratedown
migratedown:
	docker-compose run --rm api sh -c "migrate -path db/migration -database 'postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable' -verbose down"

# Run tests with coverage
.PHONY: test
test:
	docker-compose up -d
	sleep 3
	make migrateup
	@echo "🧪 Running all tests..."
	docker-compose run --rm api sh -c "go test -cover ./..."
	@echo "✅ All tests completed!"

.PHONY: test-verbose
test-verbose:
	make migrateup
	@echo "🧪 Running all tests with verbose output..."
	docker-compose run --rm api sh -c "go test -v -cover ./..."
	make migratedown

.PHONY: test-db
test-db:
	@echo "🧪 Running database tests only..."
	docker-compose run --rm api sh -c "go test -cover ./db/sqlc/"

.PHONY: test-clean
test-clean:
	@echo "🧪 Running tests with clean output..."
	./scripts/test.sh

.PHONY: test-single
test-single:
	make migrateup
	docker-compose run --rm api sh -c "go test -v -cover $(TEST) || exit 1"
	make migratedown

.PHONY: run
run:
	docker-compose up --build

.PHONY: destroy
destroy:
	docker-compose down
	make remove-data


# Remove postgres data volume
.PHONY: remove-data
remove-data:
	docker volume rm postgres_data 2>/dev/null || true

.PHONY: mockgen
mockgen:
	docker-compose run --rm api sh -c "cd /app && mockgen -package mockdb -destination db/mock/store.go lemfi/simplebank/db/sqlc Store"

.PHONY: proto
proto:
	@echo "🔧 Generating Protocol Buffer code..."
	docker-compose run --rm api sh -c "cd /app && \
		rm -f pb/*.go && \
		protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=pb --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto"
	@echo "✅ Protocol Buffer code generated!"

.PHONY: proto-clean
proto-clean:
	@echo "🧹 Cleaning Protocol Buffer generated files..."
	rm -f pb/*.go
	@echo "✅ Protocol Buffer files cleaned!"

.PHONY: evan
evan:
	evans --host localhost --port 9090 -r repl