package main

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/handlers"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/middleware/accessCheck"
	"CRM/go/apiGateway/internal/middleware/cors"
	"CRM/go/apiGateway/internal/middleware/logging"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", handlers.Authorization)
	mux.HandleFunc("/api/reg", handlers.Registration)
	mux.Handle("/api/logout", accessCheck.AccessCheck(handlers.Logout()))
	mux.HandleFunc("/api/checkAuth", handlers.CheckAuthorization)
	mux.HandleFunc("/api/updateToken", handlers.UpdateAccessToken)
	mux.Handle("/api/getUser", accessCheck.AccessCheck(handlers.GetUser()))
	mux.Handle("/api/getUser/", accessCheck.AccessCheck(handlers.GetUserById()))
	mux.Handle("/api/updateUser", accessCheck.AccessCheck(handlers.UpdateUser()))
	mux.HandleFunc("/api/getImage/", handlers.GetImage)
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world!"))
	})

	handler := logging.Logging(mux)
	handler = cors.CORS(handler)

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().ApiGateway.Address))
	var err error
	err = http.ListenAndServe(config.GetConfig().ApiGateway.Address, handler)
	/*if config.GetConfig().Env == "development" {

	} else {
		err = http.ListenAndServeTLS(config.GetConfig().ApiGateway.Address, "cert.pem", "key.pem", handler)
	}*/
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
