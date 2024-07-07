FROM golang:1.22-alpine
WORKDIR /build
COPY . .
COPY cmd/main/config.yaml .
RUN go build -o subsService cmd/main/app.go
CMD ["./subsService"]