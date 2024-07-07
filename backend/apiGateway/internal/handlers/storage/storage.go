package storage

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/proto/storageService"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func GetImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "Invalid path",
		})
		_, err := w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		return
	}

	link := parts[3]

	conn, err := grpc.DialContext(context.Background(), config.GetConfig().StorageService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for storage service",
		})
		_, err = w.Write(resp)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "connection error for storage service")
		return
	}
	defer conn.Close()

	storage := storageService.NewStorageServiceClient(conn)

	respStorage, _ := storage.GetImage(context.Background(), &storageService.GetImageRequest{
		Link: link,
	})

	if respStorage == nil || respStorage.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from storage service is nil",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "GetImage handler, response from storage service is nil")
		return
	}

	if !respStorage.Status.Successfully {
		w.WriteHeader(int(respStorage.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      respStorage.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", respStorage.Status.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(respStorage.Image)

	logger.CreateLog("info", respStorage.Status.Message)
}
