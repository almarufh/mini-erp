FROM golang:1.27.1-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd/main.go

FROM alpine:3.24.2

WORKDIR /srv

COPY --from=builder /app/main /srv/server/main

RUN mkdir -p /srv/server/public/img/

CMD [ "./server/main" ]
