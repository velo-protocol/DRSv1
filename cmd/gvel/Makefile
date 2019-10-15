test/run:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

format:
	go fmt ./...

mockgen/db:
	mockgen \
		-source=./db/interface.go \
		-destination=./db/mocks/interface.go \
		-package mocks

mockgen/friendbot:
	mockgen \
		-source=./friendbot/interface.go \
		-destination=./friendbot/mocks/interface.go \
		-package mocks