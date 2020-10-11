FROM golang:1.13 AS builder
WORKDIR /go/src/app
COPY . . 
RUN go get -d -v ./...
RUN go build -v -o GoDiscordBot

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/app/GoDiscordBot .
COPY subs.json .
CMD ["./GoDiscordBot"]
