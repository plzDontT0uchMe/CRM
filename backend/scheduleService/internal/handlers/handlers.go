package handlers

import (
	"CRM/go/scheduleService/internal/proto/scheduleService"
	"CRM/go/scheduleService/internal/service"
	"context"
)

type Server struct {
	scheduleService.UnimplementedScheduleServiceServer
}

func (s *Server) GetRecordsForUser(ctx context.Context, request *scheduleService.GetRecordsForUserRequest) (*scheduleService.GetRecordsForUserResponse, error) {
	response := &scheduleService.GetRecordsForUserResponse{}

	service.GetRecordsForUser(request, response)

	return response, nil
}

func (s *Server) GetRecords(ctx context.Context, request *scheduleService.GetRecordsRequest) (*scheduleService.GetRecordsResponse, error) {
	response := &scheduleService.GetRecordsResponse{}

	service.GetRecords(request, response)

	return response, nil
}

func (s *Server) GetRecordsByTrainerForDay(ctx context.Context, request *scheduleService.GetRecordsByTrainerForDayRequest) (*scheduleService.GetRecordsByTrainerForDayResponse, error) {
	response := &scheduleService.GetRecordsByTrainerForDayResponse{}

	service.GetRecordsByTrainerForDay(request, response)

	return response, nil
}

func (s *Server) AddRecord(ctx context.Context, request *scheduleService.AddRecordRequest) (*scheduleService.AddRecordResponse, error) {
	response := &scheduleService.AddRecordResponse{}

	service.AddRecord(request, response)

	return response, nil
}

func (s *Server) DeleteRecordById(ctx context.Context, request *scheduleService.DeleteRecordByIdRequest) (*scheduleService.DeleteRecordByIdResponse, error) {
	response := &scheduleService.DeleteRecordByIdResponse{}

	service.DeleteRecordById(request, response)

	return response, nil
}
