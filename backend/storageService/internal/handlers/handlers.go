package handlers

import (
	"CRM/go/storageService/internal/proto/storageService"
	"CRM/go/storageService/internal/service"
	"context"
)

type Server struct {
	storageService.UnimplementedStorageServiceServer
}

/*func (s *Server) Registration(ctx context.Context, registrationRequest *storageService.RegistrationRequest) (*storageService.RegistrationResponse, error) {
	err, httpStatus := service.RegisterAccount(registrationRequest)
	if err != nil {
		return &storageService.RegistrationResponse{Successfully: false, Message: "error registering account", HttpStatus: int64(httpStatus)}, nil
	}

	return &storageService.RegistrationResponse{Successfully: true, Message: "registration successfully", HttpStatus: int64(httpStatus)}, nil
}*/

func (s *Server) UploadImage(ctx context.Context, request *storageService.UploadImageRequest) (*storageService.UploadImageResponse, error) {
	response := &storageService.UploadImageResponse{}

	service.UploadImage(request, response)

	return response, nil
}

func (s *Server) GetImage(ctx context.Context, request *storageService.GetImageRequest) (*storageService.GetImageResponse, error) {
	response := &storageService.GetImageResponse{}

	service.GetImage(request, response)

	return response, nil
}

func (s *Server) GetLinkByIdAccount(ctx context.Context, request *storageService.GetLinkByIdAccountRequest) (*storageService.GetLinkByIdAccountResponse, error) {
	response := &storageService.GetLinkByIdAccountResponse{}

	service.GetLinkByIdAccount(request, response)

	return response, nil
}

func (s *Server) GetLinksByIdAccounts(ctx context.Context, request *storageService.GetLinksByIdAccountsRequest) (*storageService.GetLinksByIdAccountsResponse, error) {
	response := &storageService.GetLinksByIdAccountsResponse{}

	service.GetLinksByIdAccounts(request, response)

	return response, nil
}
