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

func UploadImage(link string, filePath string) (error, int) {
	_, err := GetDB().Exec(context.Background(), "INSERT INTO files (link, path) VALUES ($1, $2)", link, filePath)
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
