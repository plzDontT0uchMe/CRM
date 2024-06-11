package service

import (
	"CRM/go/storageService/internal/database/postgres"
	"CRM/go/storageService/internal/logger"
	"CRM/go/storageService/internal/proto/storageService"
	"CRM/go/storageService/pkg/hash"
	"CRM/go/storageService/pkg/utils"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

/*func RegisterAccount(registrationRequest *storageService.RegistrationRequest) (error, int) {
	err, httpStatus := postgres.RegisterAccount(int(registrationRequest.IdAccount))
	if err != nil {
		return err, httpStatus
	}

	return nil, httpStatus
}*/

func UploadImage(request *storageService.UploadImageRequest, response *storageService.UploadImageResponse) {
	var path sql.NullString

	row := postgres.GetPathByIdAccount(request.Id)
	err := row.Scan(&path)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan path: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error getting path",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	if path.Valid {
		_, err = postgres.DeleteImageByPath(path)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("while deleting image from database: %v", err))
			response.Status = &storageService.Status{
				Successfully: false,
				Message:      "error deleting image",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		err = os.Remove(utils.NullStringToString(path))
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("while removing image from disk: %v", err))
			response.Status = &storageService.Status{
				Successfully: false,
				Message:      "error removing image",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}
	}

	err = os.MkdirAll("img", os.ModePerm)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("while creating directory: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error creating directory",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Link = fmt.Sprintf("%v.%v", hash.GenerateRandomHash(), request.Format)
	filePath := filepath.Join("img", response.Link)

	err = os.WriteFile(filePath, request.Image, 0644)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("while writing image on disk: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error writing image",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	_, err = postgres.UploadImage(int(request.Id), response.Link, filePath)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("while uploading image in database: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error uploading image",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "upload image successfully")
	response.Status = &storageService.Status{
		Successfully: true,
		Message:      "uploading image successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetImage(request *storageService.GetImageRequest, response *storageService.GetImageResponse) {
	var path string

	row := postgres.GetPathByLink(request.Link)
	err := row.Scan(&path)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan path: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error getting path",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	image, err := os.ReadFile(path)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("while reading image: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error getting image",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Image = image

	logger.CreateLog("info", "get image successfully")
	response.Status = &storageService.Status{
		Successfully: true,
		Message:      "getting image successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetLinkByIdAccount(request *storageService.GetLinkByIdAccountRequest, response *storageService.GetLinkByIdAccountResponse) {
	var link sql.NullString

	row := postgres.GetLinkByIdAccount(request.Id)
	err := row.Scan(&link)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan link: %v", err))
		response.Status = &storageService.Status{
			Successfully: false,
			Message:      "error getting link",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Link = utils.NullStringToString(link)

	logger.CreateLog("info", "get link successfully")
	response.Status = &storageService.Status{
		Successfully: true,
		Message:      "getting link successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetLinksByIdAccounts(request *storageService.GetLinksByIdAccountsRequest, response *storageService.GetLinksByIdAccountsResponse) {
	response.Links = make(map[int64]string)
	for _, id := range request.Id {
		var link sql.NullString

		row := postgres.GetLinkByIdAccount(id)
		err := row.Scan(&link)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("scan link: %v", err))
			response.Status = &storageService.Status{
				Successfully: false,
				Message:      "error getting link",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		response.Links[id] = utils.NullStringToString(link)
	}

	logger.CreateLog("info", "get links successfully")
	response.Status = &storageService.Status{
		Successfully: true,
		Message:      "getting links successfully",
		HttpStatus:   http.StatusOK,
	}
}
