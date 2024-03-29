FROM golang:1.16-alpine AS builder
LABEL org.opencontainers.image.source=https://github.com/chand1012/Discord-Quick-Meme
WORKDIR /go/src/app
COPY . . 
RUN go get -d -v ./...
RUN go build -v -o GoDiscordBot

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/app/GoDiscordBot .
COPY subs.json .
CMD ["./GoDiscordBot"]
