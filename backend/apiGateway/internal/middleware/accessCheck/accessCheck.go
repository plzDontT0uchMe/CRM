package accessCheck

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
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

		var session authService.Session
		session.AccessToken = cookie.Value

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
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

		resp, _ := auth.CheckAuthorization(context.Background(), &authService.CheckAuthorizationRequest{
			Session: &session,
		})

		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "AuthService is crashed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "write response error")
				return
			}
			logger.CreateLog("error", "check authorization")
			return
		}

		if !resp.Status.Successfully {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      resp.Status.Message,
			})
			w.Write(res)
			logger.CreateLog("info", resp.Status.Message)
			return
		}

		logger.CreateLog("info", resp.Status.Message)
		next.ServeHTTP(w, r)
	})
}
