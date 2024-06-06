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
	account, err, httpStatus := service.AuthorizeAccount(authorizationRequest)
	if err != nil {
		return &authService.AuthorizationResponse{Successfully: false, Message: "error authorization user", HttpStatus: int64(httpStatus)}, nil
	}

	err, httpStatus = service.DeleteAllSessionsByAccount(account)
	if err != nil {
		return &authService.AuthorizationResponse{Successfully: false, Message: "error removing all sessions by user", HttpStatus: int64(httpStatus)}, nil
	}

	session, err, httpStatus := service.CreateSession(account)
	if err != nil {
		return &authService.AuthorizationResponse{Successfully: false, Message: "error creating session", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.AuthorizationResponse{Successfully: true, Message: "authorization successfully", HttpStatus: int64(httpStatus), IdAccount: int64(account.Id), Role: int64(account.Role), LastActivity: timestamppb.New(account.LastActivity), DateCreated: timestamppb.New(account.DateCreated), AccessToken: session.AccessToken, DateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken), RefreshToken: session.RefreshToken, DateExpirationRefreshToken: timestamppb.New(session.DateExpirationRefreshToken)}, nil
}

func (s *Server) Registration(ctx context.Context, registrationRequest *authService.RegistrationRequest) (*authService.RegistrationResponse, error) {
	account, err, httpStatus := service.RegisterAccount(registrationRequest)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error registering user", HttpStatus: int64(httpStatus)}, nil
	}
	session, err, httpStatus := service.CreateSession(account)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error creating session", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.RegistrationResponse{Successfully: true, Message: "registration successfully", HttpStatus: int64(httpStatus), IdAccount: int64(account.Id), Role: int64(account.Role), LastActivity: timestamppb.New(account.LastActivity), DateCreated: timestamppb.New(account.DateCreated), AccessToken: session.AccessToken, DateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken), RefreshToken: session.RefreshToken, DateExpirationRefreshToken: timestamppb.New(session.DateExpirationRefreshToken)}, nil
}

func (s *Server) Logout(ctx context.Context, logoutRequest *authService.LogoutRequest) (*authService.LogoutResponse, error) {
	err, httpStatus := service.Logout(logoutRequest)
	if err != nil {
		return &authService.LogoutResponse{Successfully: false, Message: "error logout", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.LogoutResponse{Successfully: true, Message: "logout successfully", HttpStatus: int64(httpStatus)}, nil
}

func (s *Server) CheckAuthorization(ctx context.Context, checkAuthorizationRequest *authService.CheckAuthorizationRequest) (*authService.CheckAuthorizationResponse, error) {
	account, err, httpStatus := service.CheckAuthorization(checkAuthorizationRequest)
	if err != nil {
		return &authService.CheckAuthorizationResponse{Successfully: false, Message: "authorization failed", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.CheckAuthorizationResponse{Successfully: true, Message: "authorization successfully", HttpStatus: int64(httpStatus), IdAccount: int64(account.Id), RoleAccount: int64(account.Role), LastActivityAccount: timestamppb.New(account.LastActivity), DateCreatedAccount: timestamppb.New(account.DateCreated)}, nil
}

func (s *Server) UpdateAccessToken(ctx context.Context, updateAccessTokenRequest *authService.UpdateAccessTokenRequest) (*authService.UpdateAccessTokenResponse, error) {
	session, err, httpStatus := service.UpdateAccessToken(updateAccessTokenRequest)
	if err != nil {
		return &authService.UpdateAccessTokenResponse{Successfully: false, Message: "error updating access token", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.UpdateAccessTokenResponse{Successfully: true, Message: "access token updated successfully", HttpStatus: int64(httpStatus), NewAccessToken: session.AccessToken, NewDateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken)}, nil
}

func (s *Server) GetUser(ctx context.Context, getUserRequest *authService.GetUserRequest) (*authService.GetUserResponse, error) {
	account, err, httpStatus := service.GetUser(getUserRequest)
	if err != nil {
		return &authService.GetUserResponse{Successfully: false, Message: "error getting user", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.GetUserResponse{Successfully: true, Message: "getting user successfully", HttpStatus: int64(httpStatus), Role: int64(account.Role), LastActivity: timestamppb.New(account.LastActivity), DateCreated: timestamppb.New(account.DateCreated)}, nil
}
