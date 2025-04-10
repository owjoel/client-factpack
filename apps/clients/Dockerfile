FROM golang:1.23.6-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o clients ./cmd

# ------------------------------------------

FROM alpine:latest

WORKDIR /app

RUN adduser -D -g '' app
USER app

COPY --from=build /app/clients .

EXPOSE 8081

CMD ["./clients"]
