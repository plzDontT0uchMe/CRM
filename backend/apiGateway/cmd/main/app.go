package main

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/handlers/auth"
	"CRM/go/apiGateway/internal/handlers/schedule"
	"CRM/go/apiGateway/internal/handlers/storage"
	"CRM/go/apiGateway/internal/handlers/subs"
	"CRM/go/apiGateway/internal/handlers/training"
	"CRM/go/apiGateway/internal/handlers/users"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/middleware/accessCheck"
	"CRM/go/apiGateway/internal/middleware/cors"
	"CRM/go/apiGateway/internal/middleware/logging"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", auth.Authorization)
	mux.HandleFunc("/api/reg", auth.Registration)
	mux.Handle("/api/logout", accessCheck.AccessCheck(auth.Logout()))
	mux.HandleFunc("/api/checkAuth", auth.CheckAuthorization)
	mux.HandleFunc("/api/updateToken", auth.UpdateAccessToken)

	mux.Handle("/api/getUser", accessCheck.AccessCheck(users.GetUser()))
	mux.Handle("/api/getUser/", accessCheck.AccessCheck(users.GetUserById()))
	mux.Handle("/api/getClients", accessCheck.AccessCheck(users.GetUsersByTrainer()))
	mux.Handle("/api/updateUser", accessCheck.AccessCheck(users.UpdateUser()))
	mux.HandleFunc("/api/getTrainers", users.GetTrainers)

	mux.HandleFunc("/api/getImage/", storage.GetImage)

	mux.HandleFunc("/api/getExercises", training.GetExercises)
	mux.HandleFunc("/api/getExercise/", training.GetExerciseById)
	mux.Handle("/api/createProgram", accessCheck.AccessCheck(training.CreateProgram()))
	mux.Handle("/api/getPrograms", accessCheck.AccessCheck(training.GetProgramsByUserId()))
	mux.Handle("/api/deleteProgramLocal", accessCheck.AccessCheck(training.DeleteProgramLocal()))
	mux.Handle("/api/deleteProgram", accessCheck.AccessCheck(training.DeleteProgram()))
	mux.Handle("/api/shareProgram", accessCheck.AccessCheck(training.ShareProgram()))
	mux.Handle("/api/changeProgram", accessCheck.AccessCheck(training.ChangeProgram()))

	mux.HandleFunc("/api/getSubs", subs.GetSubscriptions)
	mux.Handle("/api/createApp", accessCheck.AccessCheck(subs.CreateApplication()))
	mux.Handle("/api/getApps", accessCheck.AccessCheck(subs.GetApplications()))
	mux.Handle("/api/changeApp", accessCheck.AccessCheck(subs.ChangeSubscription()))

	mux.Handle("/api/getRecordsForUser", accessCheck.AccessCheck(schedule.GetRecordsForUser()))
	mux.Handle("/api/getRecords", accessCheck.AccessCheck(schedule.GetRecords()))
	mux.Handle("/api/addRecord", accessCheck.AccessCheck(schedule.AddRecord()))
	mux.Handle("/api/deleteRecord", accessCheck.AccessCheck(schedule.DeleteRecord()))
	mux.Handle("/api/getRecordsByTrainerForDay", accessCheck.AccessCheck(schedule.GetRecordsByTrainerForDay()))

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world!"))
	})
	mux.Handle("/api/check", accessCheck.AccessCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Check"))
	})))

	handler := logging.Logging(mux)
	handler = cors.CORS(handler)

	logger.CreateLog("info", fmt.Sprintf("starting server on %v", config.GetConfig().ApiGateway.Address))
	var err error
	err = http.ListenAndServe(config.GetConfig().ApiGateway.Address, handler)
	if err != nil {
		logger.CreateLog("error", err.Error())
	}
}
