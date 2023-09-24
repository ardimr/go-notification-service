FROM golang:1.20.4-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main cmd/app/main.go

FROM alpine:3.17 as runner
WORKDIR /run
COPY --from=builder /app/main /run/main
EXPOSE 8080
ENTRYPOINT ["./main"]