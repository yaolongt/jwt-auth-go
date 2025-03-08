FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest

RUN go mod download

EXPOSE 8000

CMD ["air"]
