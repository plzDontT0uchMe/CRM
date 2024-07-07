package handlers

import (
	"CRM/go/usersService/internal/proto/usersService"
	"CRM/go/usersService/internal/service"
	"context"
)

type Server struct {
	usersService.UsersServiceServer
}

func (s *Server) Registration(ctx context.Context, registrationRequest *usersService.RegistrationRequest) (*usersService.RegistrationResponse, error) {
	response := &usersService.RegistrationResponse{}

	service.Registration(registrationRequest, response)

	return response, nil
}

func (s *Server) GetUser(ctx context.Context, request *usersService.GetUserRequest) (*usersService.GetUserResponse, error) {
	response := &usersService.GetUserResponse{}

	service.GetUser(request, response)

	return response, nil
}

func (s *Server) GetTrainers(ctx context.Context, request *usersService.GetTrainersRequest) (*usersService.GetTrainersResponse, error) {
	response := &usersService.GetTrainersResponse{}

	service.GetTrainers(request, response)

	return response, nil
}

func (s *Server) UpdateUser(ctx context.Context, request *usersService.UpdateUserRequest) (*usersService.UpdateUserResponse, error) {
	response := &usersService.UpdateUserResponse{}

	service.UpdateUser(request, response)

	return response, nil
}
