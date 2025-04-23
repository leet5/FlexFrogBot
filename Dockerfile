FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:latest
WORKDIR /app

RUN apk add ca-certificates

COPY --from=builder /app/app .

ENV BOT_TOKEN=7541844067:AAEmnuV1Lt1v2phka_IhVib61k4Jew2v78g

ENTRYPOINT ["./app"]