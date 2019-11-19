#Docker multi-stage builds

# ------------------------------------------------------------------------------
# Development image
# ------------------------------------------------------------------------------

#Builder stage
FROM golang:1.12.7-alpine3.10 as development

# Force the go compiler to use modules
ENV GO111MODULE=on
#ENV GOPROXY=http://172.17.0.1:3333

# Update OS package and install Git
RUN apk update && apk add git mercurial bzr && apk add build-base

# Set working directory
WORKDIR /go/src/github.com/velo-protocol/DRSv1

# Install Fresh for local development
RUN go get github.com/pilu/fresh

# Install go tool for convert go test output to junit xml
RUN go get -u github.com/jstemmer/go-junit-report &&\
    go get github.com/axw/gocov/gocov &&\
    go get github.com/AlekSi/gocov-xml

# Install wait-for
RUN wget https://raw.githubusercontent.com/eficode/wait-for/master/wait-for -O /usr/local/bin/wait-for &&\
    chmod +x /usr/local/bin/wait-for

# Copy Go dependency file
COPY go.mod go.mod
COPY go.sum go.sum

# Download dependency
RUN go mod download

# Copy src
COPY node/app node/app
COPY libs libs
COPY grpc grpc

# Use CMD instead of RUN to allow command overwritability
CMD cd node/app && fresh

# ------------------------------------------------------------------------------
# Build deployable image
# ------------------------------------------------------------------------------
FROM golang:1.12.7-alpine3.10 as build

WORKDIR /go/src/github.com/velo-protocol/DRSv1

# Copy stuff from development stage
COPY --from=development /go/src/github.com/velo-protocol/DRSv1 .

# Update OS package and install Git
RUN apk update && apk add git mercurial bzr && apk add build-base

# Build the binary
ENV GO111MODULE=on
RUN cd node/app && go build -o /go/bin/node.bin

# ------------------------------------------------------------------------------
# Application image
# ------------------------------------------------------------------------------
FROM golang:1.12.7-alpine3.10

RUN apk add --no-cache tini tzdata
RUN addgroup -g 211000 -S appgroup && adduser -u 211000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

#Get artifact from builder stage
RUN mkdir -p migrations
COPY --from=build /go/bin/node.bin /app/node.bin

# Set Docker's entry point commands
RUN chown -R appuser:appgroup /app
USER appuser
ENTRYPOINT ["/sbin/tini","-sg","--","/app/node.bin"]
