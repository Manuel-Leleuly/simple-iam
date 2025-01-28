FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .  .

EXPOSE 3005

CMD [ "go", "run", "main.go" ]