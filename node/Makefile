docker/local/up:
	docker-compose -f resources/docker/docker-compose.yaml up --build

docker/local/fast:
	docker-compose -f resources/docker/docker-compose.yaml up

docker/local/app:
	docker-compose -f resources/docker/docker-compose.yaml up --build cen

docker/local/db:
	docker-compose -f resources/docker/docker-compose.yaml up -d postgres-cen-db pgadmin

docker/local/logs:
	docker-compose -f resources/docker/docker-compose.yaml logs -f

docker/local/down:
	docker-compose -f resources/docker/docker-compose.yaml down

test/run:
	export GIN_MODE=release && go test ./app/... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

test/run-debug:
	go test ./app/... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

docker/system/clean: down
	docker system prune -f
	docker volume prune -f
	docker images prune -f

format:
	go fmt ./app/...

download:
	GO111MODULE=on go mod download

tidy:
	GO111MODULE=on go mod tidy

lint:
	golangci-lint --color=auto run  --tests --exclude=comment  --exclude-use-default=false --enable=golint --enable=bodyclose --enable=gocyclo --enable=goconst --enable=unconvert ./app/...

migrate/create:
	migrate create -ext sql -dir app/migrations $(name)

proto/init:
	brew install protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go

proto/gen:
	protoc -I ./../grpc --go_out=plugins=grpc:./../grpc ./../grpc/*.proto

mockgen/repo:
	mockgen -source=./app/layers/repositories/$(m)/interface.go -destination=./app/layers/mocks/$(m)_repo.go -package mocks -mock_names Repo=Mock$(M)Repo

mockgen/usecase:
	mockgen -source=./app/layers/usecases/interface.go -destination=./app/layers/mocks/usecases.go -package mocks

mockgen/all:
	make mockgen/usecase
	make mockgen/repo m=stellar M=Stellar
	make mockgen/repo m=whitelist M=WhiteList