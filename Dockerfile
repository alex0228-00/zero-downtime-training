FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o v1 src/v1/main.go

###########################################################
# v1
###########################################################
FROM alpine:latest AS v1

WORKDIR /root/

COPY --from=builder /app/v1 .

CMD ["./v1"]
