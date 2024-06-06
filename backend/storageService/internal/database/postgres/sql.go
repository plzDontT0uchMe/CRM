package postgres

import (
	"CRM/go/storageService/internal/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func RegisterAccount(idAccount int) (error, int) {
	_, err := GetDB().Exec(context.Background(), "INSERT INTO files (id_account) VALUES ($1)", idAccount)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while registering account: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func UploadImage(idAccount int, link string, filePath string) (error, int) {
	_, err := GetDB().Exec(context.Background(), "UPDATE files SET link = $1, path = $2 WHERE id_account = $3", link, filePath, idAccount)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while uploading image: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func GetImage(link string) (*[]byte, error, int) {
	var path string
	err := GetDB().QueryRow(context.Background(), "SELECT path FROM files WHERE link = $1", link).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.CreateLog("error", fmt.Sprintf("image not found: %v", err))
			return nil, err, http.StatusBadRequest
		}
		logger.CreateLog("error", fmt.Sprintf("error while getting path: %v", err))
		return nil, err, http.StatusInternalServerError
	}

	image, err := os.ReadFile(path)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while reading file: %v", err))
		return nil, err, http.StatusInternalServerError
	}

	return &image, nil, http.StatusOK
}

func GetPathByIdAccount(idAccount int64) (*sql.NullString, error, int) {
	var path sql.NullString
	err := GetDB().QueryRow(context.Background(), "SELECT path FROM files WHERE id_account = $1", idAccount).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.CreateLog("error", fmt.Sprintf("path not found: %v", err))
			return nil, err, http.StatusBadRequest
		}
		logger.CreateLog("error", fmt.Sprintf("error while getting path: %v", err))
		return nil, err, http.StatusInternalServerError
	}

	return &path, nil, http.StatusOK
}

func GetLinkByIdAccount(idAccount int64) (*sql.NullString, error, int) {
	var link sql.NullString
	err := GetDB().QueryRow(context.Background(), "SELECT link FROM files WHERE id_account = $1", idAccount).Scan(&link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.CreateLog("error", fmt.Sprintf("link not found: %v", err))
			return nil, err, http.StatusBadRequest
		}
		logger.CreateLog("error", fmt.Sprintf("error while getting link: %v", err))
		return nil, err, http.StatusInternalServerError
	}

	return &link, nil, http.StatusOK
}

func DeleteImageByPath(path sql.NullString) (error, int) {
	_, err := GetDB().Exec(context.Background(), "UPDATE files SET link = NULL, path = NULL WHERE path = $1", path)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while deleting image: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
