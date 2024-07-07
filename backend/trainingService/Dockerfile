FROM golang:1.22-alpine
WORKDIR /build
COPY . .
COPY cmd/main/config.yaml .
RUN go build -o trainingService cmd/main/app.go
CMD ["./trainingService"]