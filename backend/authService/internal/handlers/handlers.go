package handlers

import (
	"CRM/go/authService/internal/proto/authService"
	"CRM/go/authService/internal/service"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	authService.UnimplementedAuthServiceServer
}

func (s *Server) Authorization(ctx context.Context, authorizationRequest *authService.AuthorizationRequest) (*authService.AuthorizationResponse, error) {
	user, err, httpStatus := service.AuthorizeUser(authorizationRequest)
	if err != nil {
		return &authService.AuthorizationResponse{Successfully: false, Message: "error authorization user", HttpStatus: int64(httpStatus)}, nil
	}

	err, httpStatus = service.DeleteAllSessionsByUser(user)
	if err != nil {
		return &authService.AuthorizationResponse{Successfully: false, Message: "error removing all sessions by user", HttpStatus: int64(httpStatus)}, nil
	}

	session, err, httpStatus := service.CreateSession(user)
	if err != nil {
		return &authService.AuthorizationResponse{Successfully: false, Message: "error creating session", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.AuthorizationResponse{Successfully: true, Message: "authorization successfully", HttpStatus: int64(httpStatus), AccessToken: session.AccessToken, DateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken), RefreshToken: session.RefreshToken, DateExpirationRefreshToken: timestamppb.New(session.DateExpirationRefreshToken)}, nil
}

func (s *Server) Registration(ctx context.Context, registrationRequest *authService.RegistrationRequest) (*authService.RegistrationResponse, error) {
	user, err, httpStatus := service.RegisterUser(registrationRequest)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error registering user", HttpStatus: int64(httpStatus)}, nil
	}

	err, httpStatus = service.DeleteAllSessionsByUser(user)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error removing all sessions by user", HttpStatus: int64(httpStatus)}, nil
	}

	session, err, httpStatus := service.CreateSession(user)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error creating session", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.RegistrationResponse{Successfully: true, Message: "registration successfully", HttpStatus: int64(httpStatus), AccessToken: session.AccessToken, RefreshToken: session.RefreshToken}, nil
}

func (s *Server) CheckAuthorization(ctx context.Context, checkAuthorizationRequest *authService.CheckAuthorizationRequest) (*authService.CheckAuthorizationResponse, error) {
	err, httpStatus := service.CheckAuthorization(checkAuthorizationRequest)
	if err != nil {
		return &authService.CheckAuthorizationResponse{Successfully: false, Message: "authorization failed", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.CheckAuthorizationResponse{Successfully: true, Message: "authorization successfully", HttpStatus: int64(httpStatus)}, nil
}

func (s *Server) UpdateAccessToken(ctx context.Context, updateAccessTokenRequest *authService.UpdateAccessTokenRequest) (*authService.UpdateAccessTokenResponse, error) {
	session, err, httpStatus := service.UpdateAccessToken(updateAccessTokenRequest)
	if err != nil {
		return &authService.UpdateAccessTokenResponse{Successfully: false, Message: "error updating access token", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.UpdateAccessTokenResponse{Successfully: true, Message: "access token updated successfully", HttpStatus: int64(httpStatus), NewAccessToken: session.AccessToken, NewDateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken)}, nil
}
