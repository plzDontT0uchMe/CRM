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

func (s *Server) DeleteProgramLocal(ctx context.Context, request *trainingService.DeleteProgramLocalRequest) (*trainingService.DeleteProgramLocalResponse, error) {
	response := &trainingService.DeleteProgramLocalResponse{}

	service.DeleteProgramLocal(request, response)

	return response, nil
}

func (s *Server) DeleteProgram(ctx context.Context, request *trainingService.DeleteProgramRequest) (*trainingService.DeleteProgramResponse, error) {
	response := &trainingService.DeleteProgramResponse{}

	service.DeleteProgram(request, response)

	return response, nil
}

func (s *Server) ShareProgram(ctx context.Context, request *trainingService.ShareProgramRequest) (*trainingService.ShareProgramResponse, error) {
	response := &trainingService.ShareProgramResponse{}

	service.ShareProgram(request, response)

	return response, nil
}

func (s *Server) ChangeProgram(ctx context.Context, request *trainingService.ChangeProgramRequest) (*trainingService.ChangeProgramResponse, error) {
	response := &trainingService.ChangeProgramResponse{}

	service.ChangeProgram(request, response)

	return response, nil
}
