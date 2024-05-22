package main

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/handlers"
	"CRM/go/authService/internal/logger"
	"CRM/go/authService/internal/proto/authService"
	"google.golang.org/grpc"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	authService.RegisterAuthServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().HTTPServer.Address)

	if err != nil {
		return
	}

	logger.CreateLog("info", "starting server on "+config.GetConfig().HTTPServer.Address)
	g.Serve(l)
}
