FROM golang:1.22-alpine
WORKDIR /build
COPY internal/config/config.yaml .
COPY . .
RUN go build -o apiGateway cmd/main/app.go
CMD ["./apiGateway"]