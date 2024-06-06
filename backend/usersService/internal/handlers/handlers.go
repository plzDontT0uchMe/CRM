package handlers

import (
	"CRM/go/usersService/internal/proto/usersService"
	"CRM/go/usersService/internal/service"
	"CRM/go/usersService/pkg/utils"
	"context"
)

type Server struct {
	usersService.UsersServiceServer
}

func (s *Server) Registration(ctx context.Context, registrationRequest *usersService.RegistrationRequest) (*usersService.RegistrationResponse, error) {
	err, httpStatus := service.RegisterUser(registrationRequest)
	if err != nil {
		return &usersService.RegistrationResponse{Successfully: false, Message: "error registering user", HttpStatus: int64(httpStatus)}, nil
	}

	return &usersService.RegistrationResponse{Successfully: true, Message: "registration successfully", HttpStatus: int64(httpStatus)}, nil
}

func (s *Server) GetUser(ctx context.Context, getUserRequest *usersService.GetUserRequest) (*usersService.GetUserResponse, error) {
	user, err, httpStatus := service.GetUser(getUserRequest)
	if err != nil {
		return &usersService.GetUserResponse{Successfully: false, Message: "error getting user", HttpStatus: int64(httpStatus)}, nil
	}

	return &usersService.GetUserResponse{Successfully: true, Message: "getting user successfully", HttpStatus: int64(httpStatus), Name: user.Name.String, Surname: user.Surname.String, Patronymic: user.Patronymic.String, Gender: int64(user.Gender), DateBorn: utils.ConvertSQLNullTimeToTimestamp(user.DateBorn), LinkImage: user.LinkImage.String}, nil
}

func (s *Server) UpdateUser(ctx context.Context, updateUserRequest *usersService.UpdateUserRequest) (*usersService.UpdateUserResponse, error) {
	err, httpStatus := service.UpdateUser(updateUserRequest)
	if err != nil {
		return &usersService.UpdateUserResponse{Successfully: false, Message: "error updating user", HttpStatus: int64(httpStatus)}, nil
	}

	return &usersService.UpdateUserResponse{Successfully: true, Message: "updating user successfully", HttpStatus: int64(httpStatus)}, nil
}
