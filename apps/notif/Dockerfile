FROM golang:1.24.2-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o notif ./cmd

# ------------------------------------------

FROM alpine:latest

WORKDIR /app

RUN adduser -D -g '' app
USER app

COPY --from=build /app/notif .

EXPOSE 8081

CMD ["./notif"]