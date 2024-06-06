package service

import (
	"CRM/go/storageService/internal/database/postgres"
	"CRM/go/storageService/internal/logger"
	"CRM/go/storageService/internal/proto/storageService"
	"CRM/go/storageService/pkg/hash"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func UploadImage(uploadFileRequest *storageService.UploadImageRequest) (*string, error, int) {
	err := os.MkdirAll("img", os.ModePerm)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while creating directory: %v", err))
		return nil, err, http.StatusInternalServerError
	}

	link := fmt.Sprintf("%v.%v", hash.GenerateRandomHash(), uploadFileRequest.Format)
	filePath := filepath.Join("img", link)

	err = os.WriteFile(filePath, uploadFileRequest.Image, 0644)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while saving image: %v", err))
		return nil, err, http.StatusInternalServerError
	}

	err, httpStatus := postgres.UploadImage(link, filePath)
	if err != nil {
		return nil, err, httpStatus
	}

	return &link, nil, httpStatus
}

func GetImage(getImageRequest *storageService.GetImageRequest) (*[]byte, error, int) {
	linkImage := getImageRequest.LinkImage

	image, err, httpStatus := postgres.GetImage(linkImage)
	if err != nil {
		return nil, err, httpStatus
	}

	return image, nil, httpStatus
}
