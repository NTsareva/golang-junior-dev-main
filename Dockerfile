FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o exchanges ./cmd/exchanges

FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/exchanges .
COPY configs ./configs
EXPOSE 8080
CMD ["./exchanges", "-config", "./configs/config.toml"]