FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY migrations migrations
COPY config config

RUN go build -o app cmd/app/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app
COPY migrations/versions migrations/versions
COPY config/config.yaml config/config.yaml

CMD ["./app"]