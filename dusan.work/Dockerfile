FROM golangci/golangci-lint:v1.51.1 AS lint

WORKDIR /app
COPY . .

RUN golangci-lint run -v

FROM golang:1.20.1 AS builder
COPY --from=lint /app /usr/src/app

WORKDIR /usr/src/app

RUN go mod download && go mod verify
RUN go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...
RUN go build -o dsnonl

FROM debian:buster-slim AS app
COPY --from=builder /usr/src/app/dsnonl /usr/local/bin/dsnonl

RUN mkdir /app
COPY assets /app/assets
WORKDIR /app

RUN apt-get update
RUN apt-get install -y --no-install-recommends apt-transport-https curl ca-certificates

EXPOSE 3000
CMD ["dsnonl", "server", "start"]