

protoc:
	protoc --go_out=./proto --go-grpc_out=./proto --proto_path=./proto ./proto/*.proto
	
run:
	go run cmd/main.go


docker-up:
	echo "Starting docker containers"
	docker-compose -f docker-compose.yaml up --build

docker-down:
	echo "Stopping docker containers"
	docker-compose -f docker-compose.yaml down

migrate-up:
	migrate -path migrations \
	-database "postgres://postgres:postgres@localhost:5433/auth_db?sslmode=disable" \
	up

migrate-down:
	migrate -path migrations \
	-database "postgres://postgres:postgres@localhost:5433/auth_db?sslmode=disable" \
	down 1