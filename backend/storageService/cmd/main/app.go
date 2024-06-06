package main

import (
	"CRM/go/storageService/internal/config"
	"CRM/go/storageService/internal/handlers"
	"CRM/go/storageService/internal/logger"
	"CRM/go/storageService/internal/proto/storageService"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	storageService.RegisterStorageServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().StorageService.Address)

	if err != nil {
		logger.CreateLog("error", err.Error())
		return
	}

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().StorageService.Address))
	err = g.Serve(l)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
