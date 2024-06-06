package main

import (
	"CRM/go/usersService/internal/config"
	"CRM/go/usersService/internal/handlers"
	"CRM/go/usersService/internal/logger"
	"CRM/go/usersService/internal/proto/usersService"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	usersService.RegisterUsersServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().UsersService.Address)

	if err != nil {
		logger.CreateLog("error", err.Error())
		return
	}

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().UsersService.Address))
	err = g.Serve(l)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
