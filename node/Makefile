example-docs-gen:
	@echo "============= Docs -> https://github.com/swaggo/swag ============= "
	swag init --dir ./example/docs-generator --swagger ./example/docs-generator/docs/swagger/

docs-gen:
	@echo "============= Docs -> https://github.com/swaggo/swag ============= "
	swag init --dir ./app --swagger ./app/docs/swagger/

dev:
	@echo "=============starting locally============="
	docker-compose -f resources/docker/docker-compose.yaml up --build

dev/fast:
	@echo "=============starting locally============="
	docker-compose -f resources/docker/docker-compose.yaml up

dev/app:
	@echo "=============starting locally============="
	docker-compose -f resources/docker/docker-compose.yaml up --build cen

db:
	docker-compose -f resources/docker/docker-compose.yaml up -d postgres-cen-db pgadmin

logs:
	docker-compose -f resources/docker/docker-compose.yaml logs -f

down:
	docker-compose -f resources/docker/docker-compose.yaml down

test:
	export GIN_MODE=release && go test ./app/... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

test-debug:
	go test ./app/... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

clean: down
	@echo "=============cleaning up============="
	docker system prune -f
	docker volume prune -f
	docker images prune -f

format:
	go fmt ./app/...

dep: ## Get the dependencies
	@go get -v -d ./...
	@go get -u github.com/golang/lint/golint

migrate-create:
	migrate create -ext sql -dir app/migrations $(name)

mock-module:
	mockery -dir ./app/modules/$(module) -name $(interface) -output ./app/modules/$(module)/mocks

mock-service:
	mockery -dir ./app/services/$(service) -name Interface -output ./app/services/$(service)/mocks

proto/init:
	brew install protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go

proto/gen:
	protoc -I ./../grpc --go_out=plugins=grpc:./../grpc ./../grpc/*.proto
