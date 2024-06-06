package main

import (
	"CRM/go/subsService/internal/config"
	"CRM/go/subsService/internal/logger"
	"fmt"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	authService.RegisterAuthServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().AuthService.Address)

	if err != nil {
		logger.CreateLog("error", err.Error())
		return
	}

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().AuthService.Address))
	err = g.Serve(l)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
