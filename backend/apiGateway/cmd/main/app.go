package main

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/handlers"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/middleware/accessCheck"
	"CRM/go/apiGateway/internal/middleware/cors"
	"CRM/go/apiGateway/internal/middleware/logging"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", handlers.Authorization)
	mux.HandleFunc("/api/reg", handlers.Registration)
	mux.HandleFunc("/api/checkAuth", handlers.CheckAuthorization)
	mux.HandleFunc("/api/updateToken", handlers.UpdateAccessToken)
	mux.Handle("/api/getHelloWorld", accessCheck.AccessCheck(handlers.GetHelloWorld()))

	handler := logging.Logging(mux)
	handler = cors.CORS(handler)

	logger.CreateLog("info", "starting server on "+config.GetConfig().HTTPServer.Address)
	http.ListenAndServe(config.GetConfig().HTTPServer.Address, handler)
}
