test/run:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

format:
	go fmt ./...

build:
	go build .

mockgen/db:
	mockgen \
		-source=./layers/repositories/database/interface.go \
		-destination=./layers/mocks/database.go \
		-mock_names Repository=MockDbRepo \
		-package mocks

mockgen/friendbot:
	mockgen \
		-source=./layers/repositories/friendbot/interface.go \
		-destination=./layers/mocks/friendbot.go \
		-mock_names Repository=MockFriendBotRepo \
		-package mocks

mockgen/all:
	make mockgen/db
	make mockgen/friendbot
