package handlers

import (
	"CRM/go/subsService/internal/proto/subsService"
	"CRM/go/subsService/internal/service"
	"context"
)

type Server struct {
	subsService.UnimplementedSubsServiceServer
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

func (s *Server) ChangeSubscription(ctx context.Context, request *subsService.ChangeSubscriptionRequest) (*subsService.ChangeSubscriptionResponse, error) {
	response := &subsService.ChangeSubscriptionResponse{}

	service.ChangeSubscription(request, response)

	return response, nil
}
