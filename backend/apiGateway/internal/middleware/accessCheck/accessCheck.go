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
	"strconv"
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

		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "AuthService is crashed",
			})
			w.Write(res)
			logger.CreateLog("error", "error check authorization")
			return
		}

		if !resp.Successfully {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      resp.Message,
			})
			w.Write(res)
			logger.CreateLog("info", resp.Message)
			return
		}

		logger.CreateLog("info", resp.Message)

		r.Header.Add("idAccount", strconv.FormatInt(resp.IdAccount, 10))
		r.Header.Add("roleAccount", strconv.FormatInt(resp.RoleAccount, 10))
		r.Header.Add("lastActivityAccount", resp.LastActivityAccount.AsTime().GoString())
		r.Header.Add("dateCreatedAccount", resp.DateCreatedAccount.AsTime().GoString())
		next.ServeHTTP(w, r)
	})
}
