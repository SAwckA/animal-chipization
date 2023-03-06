FROM golang:alpine AS builder

WORKDIR /build

COPY cmd cmd
COPY internal internal
COPY migrations migrations
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o app cmd/app/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app
COPY migrations/versions migrations/versions

EXPOSE 8000

CMD ["./app"]