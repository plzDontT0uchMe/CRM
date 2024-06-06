FROM golang:1.22-alpine
WORKDIR /build
COPY cmd/main/config.yaml .
COPY . .
RUN go build -o usersService cmd/main/app.go
CMD ["./usersService"]