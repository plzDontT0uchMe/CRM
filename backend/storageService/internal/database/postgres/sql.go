package postgres

import (
	"CRM/go/storageService/internal/logger"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

func RegisterAccount(idAccount int) (error, int) {
	_, err := GetDB().Exec(context.Background(), "INSERT INTO files (id_account) VALUES ($1)", idAccount)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error while registering account: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func UploadImage(idAccount int, link string, filePath string) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE files SET link = $1, path = $2 WHERE id_account = $3", link, filePath, idAccount)
}

func GetPathByLink(link string) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT path FROM files WHERE link = $1", link)
}

func GetPathByIdAccount(idAccount int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT path FROM files WHERE id_account = $1", idAccount)
}

func GetLinkByIdAccount(idAccount int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT link FROM files WHERE id_account = $1", idAccount)
}

func DeleteImageByPath(path sql.NullString) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE files SET link = NULL, path = NULL WHERE path = $1", path)
}
