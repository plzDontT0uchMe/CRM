package handlers

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/proto/authService"
	"CRM/go/authService/internal/proto/usersService"
	"CRM/go/authService/internal/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
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

	return &authService.AuthorizationResponse{Successfully: true, Message: "authorization successfully", HttpStatus: int64(httpStatus), AccessToken: session.AccessToken, DateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken), RefreshToken: session.RefreshToken, DateExpirationRefreshToken: timestamppb.New(session.DateExpirationRefreshToken)}, nil
}

func (s *Server) Registration(ctx context.Context, registrationRequest *authService.RegistrationRequest) (*authService.RegistrationResponse, error) {
	account, err, httpStatus := service.RegisterAccount(registrationRequest)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error registering user", HttpStatus: int64(httpStatus)}, nil
	}

	conn, err := grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "connection error for users service", HttpStatus: http.StatusInternalServerError}, nil
	}
	defer conn.Close()

	users := usersService.NewUsersServiceClient(conn)

	resp, _ := users.Registration(context.Background(), &usersService.RegistrationRequest{IdAccount: int64(account.Id)})

	if resp == nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "response from users service is nil", HttpStatus: http.StatusInternalServerError}, nil
	}
	if !resp.Successfully {
		return &authService.RegistrationResponse{Successfully: false, Message: resp.Message, HttpStatus: resp.HttpStatus}, nil
	}

	session, err, httpStatus := service.CreateSession(account)
	if err != nil {
		return &authService.RegistrationResponse{Successfully: false, Message: "error creating session", HttpStatus: int64(httpStatus)}, nil
	}

	return &authService.RegistrationResponse{Successfully: true, Message: "registration successfully", HttpStatus: int64(httpStatus), AccessToken: session.AccessToken, DateExpirationAccessToken: timestamppb.New(session.DateExpirationAccessToken), RefreshToken: session.RefreshToken, DateExpirationRefreshToken: timestamppb.New(session.DateExpirationRefreshToken)}, nil
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
