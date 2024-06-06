package service

import (
	"CRM/go/storageService/internal/database/postgres"
	"CRM/go/storageService/internal/logger"
	"CRM/go/storageService/internal/proto/storageService"
	"CRM/go/storageService/pkg/hash"
	"CRM/go/storageService/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterAccount(registrationRequest *storageService.RegistrationRequest) (error, int) {
	err, httpStatus := postgres.RegisterAccount(int(registrationRequest.IdAccount))
	if err != nil {
		return err, httpStatus
	}

	return nil, httpStatus
}

func UploadImage(uploadFileRequest *storageService.UploadImageRequest) (*string, error, int) {
	path, err, httpStatus := postgres.GetPathByIdAccount(uploadFileRequest.IdAccount)
	if err == nil && path.Valid {
		err, httpStatus = postgres.DeleteImageByPath(*path)
		if err != nil {
			return nil, err, httpStatus
		}
		err = os.Remove(utils.NullStringToString(*path))
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("error while removing image: %v", err))
			return nil, err, http.StatusInternalServerError
		}
	}

	err = os.MkdirAll("img", os.ModePerm)
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

	err, httpStatus = postgres.UploadImage(int(uploadFileRequest.IdAccount), link, filePath)
	if err != nil {
		return nil, err, httpStatus
	}

	return &link, nil, httpStatus
}

func GetImage(getImageRequest *storageService.GetImageRequest) (*[]byte, error, int) {
	linkImage := getImageRequest.Link

	image, err, httpStatus := postgres.GetImage(linkImage)
	if err != nil {
		return nil, err, httpStatus
	}

	return image, nil, httpStatus
}

func GetLinkByIdAccount(getLinkByIdAccountRequest *storageService.GetLinkByIdAccountRequest) (*string, error, int) {
	link, err, httpStatus := postgres.GetLinkByIdAccount(getLinkByIdAccountRequest.IdAccount)
	if err != nil {
		return nil, err, httpStatus
	}

	linkString := utils.NullStringToString(*link)

	return &linkString, nil, httpStatus
}
