FROM golang:1.22-alpine
WORKDIR /build
COPY cmd/main/config.yaml .
COPY . .
RUN go build -o storageService cmd/main/app.go
CMD ["./storageService"]