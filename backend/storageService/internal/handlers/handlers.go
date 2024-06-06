package handlers

import (
	"CRM/go/storageService/internal/proto/storageService"
	"CRM/go/storageService/internal/service"
	"context"
)

type Server struct {
	storageService.UnimplementedStorageServiceServer
}

func (s *Server) UploadImage(ctx context.Context, uploadFileRequest *storageService.UploadImageRequest) (*storageService.UploadImageResponse, error) {
	link, err, httpStatus := service.UploadImage(uploadFileRequest)
	if err != nil {
		return &storageService.UploadImageResponse{Successfully: false, Message: "error uploading image", HttpStatus: int64(httpStatus)}, nil
	}

	return &storageService.UploadImageResponse{Successfully: true, Message: "uploading image successfully", HttpStatus: int64(httpStatus), Link: *link}, nil
}

func (s *Server) GetImage(ctx context.Context, getImageRequest *storageService.GetImageRequest) (*storageService.GetImageResponse, error) {
	image, err, httpStatus := service.GetImage(getImageRequest)
	if err != nil {
		return &storageService.GetImageResponse{Successfully: false, Message: "error getting image", HttpStatus: int64(httpStatus)}, nil
	}

	return &storageService.GetImageResponse{Successfully: true, Message: "getting image successfully", HttpStatus: int64(httpStatus), Image: *image}, nil
}
