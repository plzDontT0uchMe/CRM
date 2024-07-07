package service

import (
	"CRM/go/subsService/internal/database/postgres"
	"CRM/go/subsService/internal/logger"
	"CRM/go/subsService/internal/proto/subsService"
	"CRM/go/subsService/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strings"
	"time"
)

func Registration(request *subsService.RegistrationRequest, response *subsService.RegistrationResponse) {
	rows, err := postgres.GetSubscriptions()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting subscriptions",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	subscription := &subsService.Subscription{}

	for rows.Next() {
		var possibilities string

		err = rows.Scan(&subscription.Id, &subscription.Name, &subscription.Price, &subscription.Description, &possibilities)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
			response.Status = &subsService.Status{
				Successfully: false,
				Message:      "error getting subscriptions",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		if subscription.Price == 0 {
			break
		}
	}

	_, err = postgres.Registration(request.Id, subscription.Id)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error registering account",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "account successfully registered")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "account successfully registered",
		HttpStatus:   http.StatusOK,
	}
}

func GetSubscriptions(request *subsService.GetSubscriptionsRequest, response *subsService.GetSubscriptionsResponse) {
	rows, err := postgres.GetSubscriptions()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting subscriptions",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		subscription := &subsService.Subscription{}
		var possibilities string

		err = rows.Scan(&subscription.Id, &subscription.Name, &subscription.Price, &subscription.Description, &possibilities)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
			response.Status = &subsService.Status{
				Successfully: false,
				Message:      "error getting subscriptions",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		subscription.Possibilities = strings.Split(possibilities, ",")
		response.Subscriptions = append(response.Subscriptions, subscription)
	}

	logger.CreateLog("info", "subscriptions successfully received")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "subscriptions successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func GetSubscriptionByAccountId(request *subsService.GetSubscriptionByAccountIdRequest, response *subsService.GetSubscriptionByAccountIdResponse) {
	response.Subscription = &subsService.Subscription{}

	var possibilities string
	var dateExpiration sql.NullTime
	var idTrainer sql.NullInt64

	row := postgres.GetSubscriptionByAccountId(request.Id)
	err := row.Scan(&response.Subscription.Id, &response.Subscription.Name, &response.Subscription.Price, &response.Subscription.Description, &possibilities, &idTrainer, &dateExpiration)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting subscription",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Subscription.Possibilities = strings.Split(possibilities, ",")
	response.Subscription.DateExpiration = utils.ConvertSQLNullTimeToTimestamp(dateExpiration)
	response.Subscription.IdTrainer = idTrainer.Int64

	logger.CreateLog("info", "subscription successfully received")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "subscription successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func GetSubscriptionById(request *subsService.GetSubscriptionByIdRequest, response *subsService.GetSubscriptionByIdResponse) {
	response.Subscription = &subsService.Subscription{}

	var possibilities string

	row := postgres.GetSubscriptionById(request.Id)
	err := row.Scan(&response.Subscription.Id, &response.Subscription.Name, &response.Subscription.Price, &response.Subscription.Description, &possibilities)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting subscription",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Subscription.Possibilities = strings.Split(possibilities, ",")

	logger.CreateLog("info", "subscription successfully received")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "subscription successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func ChangeApplication(request *subsService.ChangeApplicationRequest, response *subsService.ChangeApplicationResponse) {
	if request.IsAccepted {
		application := subsService.Application{
			Subscription: &subsService.Subscription{},
		}

		idTrainer := sql.NullInt64{}
		var possibilities string
		var duration sql.NullInt64

		row := postgres.GetApplicationById(request.Id)
		err := row.Scan(&application.Id, &application.IdClient, &idTrainer, &application.Subscription.Id, &application.Subscription.Name, &application.Subscription.Price, &application.Subscription.Description, &possibilities, &duration)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
			response.Status = &subsService.Status{
				Successfully: false,
				Message:      "error changing subscription",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		if idTrainer.Int64 == 0 {
			idTrainer.Valid = false
		}

		var dateExpiration sql.NullTime
		if duration.Int64 != 0 {
			dateExpiration = utils.ConvertTimestampToSQLNullTime(timestamppb.New(time.Now().Add(time.Duration(duration.Int64) * 30 * 24 * time.Hour)))
		} else {
			dateExpiration.Valid = false
		}

		_, err = postgres.ChangeSubscription(application.IdClient, application.Subscription.Id, idTrainer, dateExpiration)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
			response.Status = &subsService.Status{
				Successfully: false,
				Message:      "error changing subscription",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}
	}

	_, err := postgres.DeleteApplicationById(request.Id)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error changing subscription",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "subscription successfully changed")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "subscription successfully changed",
		HttpStatus:   http.StatusOK,
	}
}

func CreateApplication(request *subsService.CreateApplicationRequest, response *subsService.CreateApplicationResponse) {
	idTrainer := utils.ConvertInt64ToSQLNullInt64(request.Application.IdTrainer)
	if request.Application.IdTrainer == 0 {
		idTrainer.Valid = false
	}
	duration := utils.ConvertInt64ToSQLNullInt64(request.Application.Duration)
	if request.Application.Duration == 0 {
		duration.Valid = false
	}

	_, err := postgres.CreateApplication(request.Application.IdClient, request.Application.Subscription.Id, idTrainer, duration)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error creating application",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "application successfully created")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "application successfully created",
		HttpStatus:   http.StatusOK,
	}
}

func GetApplications(request *subsService.GetApplicationsRequest, response *subsService.GetApplicationsResponse) {
	response.Applications = []*subsService.Application{}

	rows, err := postgres.GetApplications()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting applications",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		application := &subsService.Application{
			Subscription: &subsService.Subscription{},
		}
		var possibilities string
		var duration sql.NullInt64
		var idTrainer sql.NullInt64

		err = rows.Scan(&application.Id, &application.IdClient, &idTrainer, &application.Subscription.Id, &application.Subscription.Name, &application.Subscription.Price, &application.Subscription.Description, &possibilities, &duration)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
			response.Status = &subsService.Status{
				Successfully: false,
				Message:      "error getting applications",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		application.IdTrainer = idTrainer.Int64
		application.Subscription.Possibilities = strings.Split(possibilities, ",")
		application.Duration = duration.Int64

		response.Applications = append(response.Applications, application)
	}

	logger.CreateLog("info", "applications successfully received")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "applications successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func GetSubscriptionAndApplicationByAccountId(request *subsService.GetSubscriptionAndApplicationByAccountIdRequest, response *subsService.GetSubscriptionAndApplicationByAccountIdResponse) {
	response.Subscription = &subsService.Subscription{}
	response.Application = &subsService.Application{
		Subscription: &subsService.Subscription{},
	}

	var possibilitiesForSub string
	var dateExpiration sql.NullTime
	var idTrainerForSub sql.NullInt64

	var possibilitiesForApp string
	var duration sql.NullInt64
	var idTrainerForApp sql.NullInt64

	row := postgres.GetSubscriptionByAccountId(request.Id)
	err := row.Scan(&response.Subscription.Id, &response.Subscription.Name, &response.Subscription.Price, &response.Subscription.Description, &possibilitiesForSub, &idTrainerForSub, &dateExpiration)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting subscription and application",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Subscription.IdTrainer = idTrainerForSub.Int64
	response.Subscription.Possibilities = strings.Split(possibilitiesForSub, ",")
	response.Subscription.DateExpiration = utils.ConvertSQLNullTimeToTimestamp(dateExpiration)

	row = postgres.GetApplicationByAccountId(request.Id)
	err = row.Scan(&response.Application.Id, &response.Application.IdClient, &idTrainerForApp, &response.Application.Subscription.Id, &response.Application.Subscription.Name, &response.Application.Subscription.Price, &response.Application.Subscription.Description, &possibilitiesForApp, &duration)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			logger.CreateLog("info", "subscription and application successfully received")
			response.Status = &subsService.Status{
				Successfully: true,
				Message:      "subscription successfully received",
				HttpStatus:   http.StatusOK,
			}
			return
		}
		logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting subscription and application",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	response.Application.IdTrainer = idTrainerForApp.Int64
	response.Application.Subscription.Possibilities = strings.Split(possibilitiesForApp, ",")
	response.Application.Duration = duration.Int64

	logger.CreateLog("info", "subscription and application successfully received")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "subscription and application successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func GetUsersByTrainerId(request *subsService.GetUsersByTrainerIdRequest, response *subsService.GetUsersByTrainerIdResponse) {
	rows, err := postgres.GetUsersByTrainerId(request.Id)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &subsService.Status{
			Successfully: false,
			Message:      "error getting users",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user sql.NullInt64
		err = rows.Scan(&user)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
			response.Status = &subsService.Status{
				Successfully: false,
				Message:      "error getting users",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		response.Id = append(response.Id, user.Int64)
	}

	logger.CreateLog("info", "users successfully received")
	response.Status = &subsService.Status{
		Successfully: true,
		Message:      "users successfully received",
		HttpStatus:   http.StatusOK,
	}
}
