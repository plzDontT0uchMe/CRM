package handlers

import (
	"CRM/go/trainingService/internal/proto/trainingService"
	"CRM/go/trainingService/internal/service"
	"context"
)

type Server struct {
	trainingService.UnimplementedTrainingServiceServer
}

func (s *Server) GetExercises(ctx context.Context, request *trainingService.GetExercisesRequest) (*trainingService.GetExercisesResponse, error) {
	response := &trainingService.GetExercisesResponse{}

	service.GetExercises(request, response)

	return response, nil
}

func (s *Server) GetExerciseById(ctx context.Context, request *trainingService.GetExerciseByIdRequest) (*trainingService.GetExerciseByIdResponse, error) {
	response := &trainingService.GetExerciseByIdResponse{}

	service.GetExerciseById(request, response)

	return response, nil
}

func (s *Server) CreateProgram(ctx context.Context, request *trainingService.CreateProgramRequest) (*trainingService.CreateProgramResponse, error) {
	response := &trainingService.CreateProgramResponse{}

	service.CreateProgram(request, response)

	return response, nil
}

func (s *Server) GetProgramsByUserId(ctx context.Context, request *trainingService.GetProgramsByUserIdRequest) (*trainingService.GetProgramsByUserIdResponse, error) {
	response := &trainingService.GetProgramsByUserIdResponse{}

	service.GetProgramsByUserId(request, response)

	return response, nil
}
