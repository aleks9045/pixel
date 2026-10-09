FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/app/main.go

FROM alpine

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/.env .

EXPOSE 8080

CMD ["./app"]