FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o translate_api cmd/translate_api/translate_api.go

# Remove code compiler and tool chain from the deployed image
FROM alpine:3.20

COPY .env .env

COPY --from=builder /app/translate_api /translate_api

ENV PORT=8080

EXPOSE $PORT

# Command to run the executable
CMD ["/translate_api"]
