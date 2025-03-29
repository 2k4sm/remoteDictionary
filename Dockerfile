FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o remoteDictionary

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/remoteDictionary .

EXPOSE 7171

CMD ["./remoteDictionary"]
