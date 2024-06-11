package main

import (
	"CRM/go/trainingService/internal/config"
	"CRM/go/trainingService/internal/handlers"
	"CRM/go/trainingService/internal/logger"
	"CRM/go/trainingService/internal/proto/trainingService"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	trainingService.RegisterTrainingServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().TrainingService.Address)

	if err != nil {
		logger.CreateLog("error", err.Error())
		return
	}

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().TrainingService.Address))
	err = g.Serve(l)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
