package service

import (
	"CRM/go/subsService/internal/database/postgres"
	"CRM/go/subsService/internal/logger"
	"CRM/go/subsService/internal/proto/subsService"
	"CRM/go/subsService/pkg/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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

	fmt.Println(request.Id)
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
	fmt.Println(response.Subscription)

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

func ChangeSubscription(request *subsService.ChangeSubscriptionRequest, response *subsService.ChangeSubscriptionResponse) {
	dateExpiration := utils.ConvertTimestampToSQLNullTime(request.DateExpiration)
	if dateExpiration.Time.Before(time.Now().UTC()) {
		dateExpiration.Valid = false
	}

	idTrainer := utils.ConvertInt64ToSQLNullInt64(request.IdTrainer)
	if request.IdTrainer == 0 {
		idTrainer.Valid = false
	}

	_, err := postgres.ChangeSubscription(request.IdClient, request.IdSubscription, idTrainer, dateExpiration)
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
