package handlers

import (
	"CRM/go/authService/internal/proto/authService"
	"CRM/go/authService/internal/service"
	"golang.org/x/net/context"
)

type Server struct {
	authService.UnimplementedAuthServiceServer
}

func (s *Server) Authorization(ctx context.Context, request *authService.AuthorizationRequest) (*authService.AuthorizationResponse, error) {
	response := &authService.AuthorizationResponse{}

	service.Authorization(request, response)

	return response, nil
}

func (s *Server) Registration(ctx context.Context, request *authService.RegistrationRequest) (*authService.RegistrationResponse, error) {
	response := &authService.RegistrationResponse{}

	service.CreateAccount(request, response)

	return response, nil
}

func (s *Server) Logout(ctx context.Context, request *authService.LogoutRequest) (*authService.LogoutResponse, error) {
	response := &authService.LogoutResponse{}

	service.Logout(request, response)

	return response, nil
}

func (s *Server) CheckAuthorization(ctx context.Context, request *authService.CheckAuthorizationRequest) (*authService.CheckAuthorizationResponse, error) {
	response := &authService.CheckAuthorizationResponse{}

	service.CheckAuthorization(request, response)

	return response, nil
}

func (s *Server) UpdateAccessToken(ctx context.Context, request *authService.UpdateAccessTokenRequest) (*authService.UpdateAccessTokenResponse, error) {
	response := &authService.UpdateAccessTokenResponse{}

	service.UpdateAccessToken(request, response)

	return response, nil
}

func (s *Server) GetAccountByAccessToken(ctx context.Context, request *authService.GetAccountByAccessTokenRequest) (*authService.GetAccountByAccessTokenResponse, error) {
	response := &authService.GetAccountByAccessTokenResponse{}

	service.GetAccount(request, response)

	return response, nil
}

func (s *Server) GetAccountById(ctx context.Context, request *authService.GetAccountByIdRequest) (*authService.GetAccountByIdResponse, error) {
	response := &authService.GetAccountByIdResponse{}

	service.GetAccountById(request, response)

	return response, nil
}

func (s *Server) GetAccounts(ctx context.Context, request *authService.GetAccountsRequest) (*authService.GetAccountsResponse, error) {
	response := &authService.GetAccountsResponse{}

	service.GetAccounts(request, response)

	return response, nil
}
