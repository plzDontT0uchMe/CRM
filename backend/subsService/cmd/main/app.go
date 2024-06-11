package main

import (
	"CRM/go/subsService/internal/config"
	"CRM/go/subsService/internal/handlers"
	"CRM/go/subsService/internal/logger"
	"CRM/go/subsService/internal/proto/subsService"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	subsService.RegisterSubsServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().SubsService.Address)

	if err != nil {
		logger.CreateLog("error", err.Error())
		return
	}

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().SubsService.Address))
	err = g.Serve(l)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
