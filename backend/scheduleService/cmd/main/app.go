package main

import (
	"CRM/go/scheduleService/internal/config"
	"CRM/go/scheduleService/internal/handlers"
	"CRM/go/scheduleService/internal/logger"
	"CRM/go/scheduleService/internal/proto/scheduleService"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	g := grpc.NewServer()

	srv := &handlers.Server{}

	scheduleService.RegisterScheduleServiceServer(g, srv)

	l, err := net.Listen("tcp", config.GetConfig().ScheduleService.Address)

	if err != nil {
		logger.CreateLog("error", err.Error())
		return
	}

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().ScheduleService.Address))
	err = g.Serve(l)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
