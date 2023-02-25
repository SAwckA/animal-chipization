FROM golang:alpine AS builder

WORKDIR /build

COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o app cmd/app/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app

EXPOSE 8000

CMD ["./app"]