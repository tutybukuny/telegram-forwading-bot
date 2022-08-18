FROM golang:1.18-alpine as builder
RUN apk add --no-cache --update gcc musl-dev g++ make git gnutls gnutls-dev gnutls-c++ bash git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd cmd
COPY config config
COPY pkg pkg
COPY internal internal

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags musl -o ./dist/bot ./cmd/bot


FROM alpine:3.15.4 AS target
WORKDIR /app

RUN apk add gettext libintl && rm -rf /var/cache/apk/*

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dist/bot /app/bot

CMD ["/app/bot"]