# Stage 1: build with CGO disabled for static binary
FROM golang:1.23.5 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY cmd/mockclient ./cmd/mockclient
RUN go build -o transit-api ./cmd/mockclient

# Stage 2: minimal runtime
FROM scratch

COPY --from=builder /app/transit-api /transit-api

EXPOSE 8080
ENTRYPOINT ["/transit-api"]
