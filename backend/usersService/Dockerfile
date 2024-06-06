FROM golang:1.22-alpine
WORKDIR /build
COPY internal/config/config.yaml .
COPY . .
RUN go build -o authService cmd/main/app.go
CMD ["./authService"]