FROM golang:1.21-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CG0_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]