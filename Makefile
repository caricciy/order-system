DB_PORT ?= 3307
DB_USER ?= root
DB_PASSWORD ?= root
PROTO_DIR ?=./internal/infra/grpc/protofiles
OUT_DIR ?=./internal/infra/grpc/pb
PROTO_FILES ?=$(wildcard $(PROTO_DIR)/*.proto)
RABBITMQ_CONTAINER ?=rabbitmq-3
RABBITMQ_USER ?= guest
RABBITMQ_PASSWORD ?= guest

create-migration:
	@echo "Creating a new migration"
	migrate create -ext=sql -dir=sql/migrations -seq $(name)

migrate:
	@echo "Applying migrations"
	migrate -path=internal/sql/migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp(localhost:$(DB_PORT))/orders" -verbose up

rollback:
	@echo "Rolling back the last migration"
	migrate -path=internal/sql/migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp(localhost:$(DB_PORT))/orders" -verbose down

## Generate GraphQL code using gqlgen
graphql-gen:
	@echo "Generating GraphQL code"
	go run github.com/99designs/gqlgen generate

## Generate gRPC code from .proto files
grpc-gen:
	@echo "Generating gRPC code"
	@mkdir -p $(OUT_DIR)
	@protoc -I=$(PROTO_DIR) \
       --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
       $(PROTO_FILES)
	@go mod tidy

## Clean up generated gRPC files
grpc-clean:
	@rm -rf $(OUT_DIR)/*.pb.go

## Run Evans REPL for gRPC service
evans-repl:
	evans -r repl -p 50051 --package pb --service OrderService

test:
	@echo "Running tests"
	go test ./... -v

## Run the server
server-start: rabbitmq-queue migrate
	@clear
	@echo "Starting the server"
	go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

## This is a helper task to create a RabbitMQ queue and bind it to the default exchange
rabbitmq-queue:
	@docker exec -it $(RABBITMQ_CONTAINER) rabbitmqadmin --username=$(RABBITMQ_USER) --password=$(RABBITMQ_PASSWORD) declare queue name="order.created" durable=true
	@clear
	@docker exec -it $(RABBITMQ_CONTAINER) rabbitmqadmin --username=$(RABBITMQ_USER) --password=$(RABBITMQ_PASSWORD) declare binding source="amq.direct" destination_type="queue" destination="order.created"
	@clear
.PHONY: migrate rollback create-migration graphql-gen grpc-gen grpc-clean test server-start