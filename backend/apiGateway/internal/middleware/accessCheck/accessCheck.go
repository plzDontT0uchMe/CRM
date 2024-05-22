package accessCheck

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"net/http"
)

func AccessCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			w.Write(resp)
			logger.CreateLog("error", "CheckAuthorization handler, access token not found in cookie")
			return
		}

		var session models.Session
		session.AccessToken = cookie.Value

		conn, err := grpc.Dial(config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for auth service",
			})
			w.Write(resp)
			logger.CreateLog("error", "connection error for auth service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		resp, _ := auth.CheckAuthorization(context.Background(), &authService.CheckAuthorizationRequest{AccessToken: session.AccessToken})

		logger.CreateLog("info", resp.Message)

		if resp.Successfully {
			next.ServeHTTP(w, r)
		}
	})
}
