FROM golang:1.23.6-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth ./cmd

# ------------------------------------------

FROM alpine:3.18

WORKDIR /app

RUN adduser -D -g '' app
USER app

COPY --from=build /app/auth .

EXPOSE 8080

CMD ["./auth"]