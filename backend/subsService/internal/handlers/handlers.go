package handlers

import (
	"CRM/go/subsService/internal/proto/subsService"
	"CRM/go/subsService/internal/service"
	"context"
)

type Server struct {
	subsService.UnimplementedSubsServiceServer
}

func (s *Server) Registration(ctx context.Context, request *subsService.RegistrationRequest) (*subsService.RegistrationResponse, error) {
	response := &subsService.RegistrationResponse{}

	service.Registration(request, response)

	return response, nil
}

func (s *Server) GetSubscriptions(ctx context.Context, request *subsService.GetSubscriptionsRequest) (*subsService.GetSubscriptionsResponse, error) {
	response := &subsService.GetSubscriptionsResponse{}

	service.GetSubscriptions(request, response)

	return response, nil
}

func (s *Server) GetSubscriptionByAccountId(ctx context.Context, request *subsService.GetSubscriptionByAccountIdRequest) (*subsService.GetSubscriptionByAccountIdResponse, error) {
	response := &subsService.GetSubscriptionByAccountIdResponse{}

	service.GetSubscriptionByAccountId(request, response)

	return response, nil
}

func (s *Server) GetSubscriptionById(ctx context.Context, request *subsService.GetSubscriptionByIdRequest) (*subsService.GetSubscriptionByIdResponse, error) {
	response := &subsService.GetSubscriptionByIdResponse{}

	service.GetSubscriptionById(request, response)

	return response, nil
}

func (s *Server) ChangeApplication(ctx context.Context, request *subsService.ChangeApplicationRequest) (*subsService.ChangeApplicationResponse, error) {
	response := &subsService.ChangeApplicationResponse{}

	service.ChangeApplication(request, response)

	return response, nil
}

func (s *Server) CreateApplication(ctx context.Context, request *subsService.CreateApplicationRequest) (*subsService.CreateApplicationResponse, error) {
	response := &subsService.CreateApplicationResponse{}

	service.CreateApplication(request, response)

	return response, nil
}

func (s *Server) GetApplications(ctx context.Context, request *subsService.GetApplicationsRequest) (*subsService.GetApplicationsResponse, error) {
	response := &subsService.GetApplicationsResponse{}

	service.GetApplications(request, response)

	return response, nil
}

func (s *Server) GetSubscriptionAndApplicationByAccountId(ctx context.Context, request *subsService.GetSubscriptionAndApplicationByAccountIdRequest) (*subsService.GetSubscriptionAndApplicationByAccountIdResponse, error) {
	response := &subsService.GetSubscriptionAndApplicationByAccountIdResponse{}

	service.GetSubscriptionAndApplicationByAccountId(request, response)

	return response, nil
}

func (s *Server) GetUsersByTrainerId(ctx context.Context, request *subsService.GetUsersByTrainerIdRequest) (*subsService.GetUsersByTrainerIdResponse, error) {
	response := &subsService.GetUsersByTrainerIdResponse{}

	service.GetUsersByTrainerId(request, response)

	return response, nil
}
