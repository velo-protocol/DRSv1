#Docker multi-stage builds

# ------------------------------------------------------------------------------
# Development image
# ------------------------------------------------------------------------------

#Builder stage
FROM golang:1.12.7-alpine3.10 as builder
# Update OS package and install Git
RUN apk update && apk add git && apk add build-base
# Set working directory
WORKDIR /go/src/gitlab.com/velo-labs/cen
# Install Fresh for local development
RUN go get github.com/pilu/fresh
# Install go tool for convert go test output to junit xml
RUN go get -u github.com/jstemmer/go-junit-report &&\
    go get github.com/axw/gocov/... &&\
    go get github.com/AlekSi/gocov-xml &&\
    go get -u golang.org/x/lint/golint
# Install wait-for
RUN wget https://raw.githubusercontent.com/eficode/wait-for/master/wait-for -O /usr/local/bin/wait-for &&\
    chmod +x /usr/local/bin/wait-for
# Copy Go dependency file
COPY go.mod go.mod
COPY go.sum go.sum
# Run go mod tidy
RUN go mod download
# Copy Go source code
COPY . .
# Set Docker's entry point commands
RUN cd node/app/ && go build

# ------------------------------------------------------------------------------
# Deployment image
# ------------------------------------------------------------------------------

#App stage
FROM golang:1.12.7-alpine3.10
# Set working directory
WORKDIR /go/src/gitlab.com/velo-labs/cen
#Get artifact from buiber stage
COPY --from=builder /usr/local/bin/wait-for /usr/local/bin/wait-for
COPY --from=builder /go/src/gitlab.com/velo-labs/cen/node/app app/
COPY --from=builder /go/src/gitlab.com/velo-labs/cen/vendor vendor
# Set Docker's entry point commands
CMD cd app/ && ./app;
