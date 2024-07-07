FROM golang:1.22-alpine
WORKDIR /build
COPY . .
COPY cmd/main/config.yaml .
COPY cmd/main/img /build/img
COPY cmd/main/exercise /build/exercise
RUN go build -o storageService cmd/main/app.go
CMD ["./storageService"]