test/run:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

format:
	go fmt ./...

gvel/build:
	go build .

mockgen/logic:
	mockgen \
		-source=./layers/logic/interface.go \
		-destination=./layers/mocks/logic.go \
		-mock_names Logic=MockLogic \
		-package mocks

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

mockgen/console:
	mockgen \
		-source=./utils/console/interface.go \
		-destination=./utils/mocks/console.go \
		-mock_names Prompt=MockPrompt \
		-package mockutils

mockgen/config:
	mockgen \
		-source=./utils/config/interface.go \
		-destination=./utils/mocks/configuration.go \
		-mock_names Prompt=MockConfiguration \
		-package mockutils

mockgen/all:
	make mockgen/db
	make mockgen/friendbot
	make mockgen/logic
	make mockgen/console
	make mockgen/config
