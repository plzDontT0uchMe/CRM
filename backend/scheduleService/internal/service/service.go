package service

import (
	"CRM/go/scheduleService/internal/database/postgres"
	"CRM/go/scheduleService/internal/logger"
	"CRM/go/scheduleService/internal/proto/scheduleService"
	"CRM/go/scheduleService/pkg/utils"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetRecordsForUser(request *scheduleService.GetRecordsForUserRequest, response *scheduleService.GetRecordsForUserResponse) {
	records, err := postgres.GetRecordsForUser(request.UserId)
	if err != nil {
		logger.CreateLog("error", "get records for user: %v", err)
		response.Status = &scheduleService.Status{
			Successfully: false,
			Message:      "error getting records for user",
		}
		return
	}

	for records.Next() {
		var trainerId sql.NullInt64
		var dateStart, dateEnd sql.NullTime

		record := &scheduleService.Record{}
		err = records.Scan(&record.Id, &record.ClientId, &trainerId, &dateStart, &dateEnd)
		if err != nil {
			logger.CreateLog("error", "scan records for user: %v", err)
			response.Status = &scheduleService.Status{
				Successfully: false,
				Message:      "error scanning records for user",
			}
			return
		}

		if trainerId.Valid {
			record.TrainerId = trainerId.Int64
		}

		record.DateStart = timestamppb.New(dateStart.Time)
		record.DateEnd = timestamppb.New(dateEnd.Time)

		response.Records = append(response.Records, record)

	}
	response.Status = &scheduleService.Status{
		Successfully: true,
		Message:      "records for user received successfully",
	}
	logger.CreateLog("info", "records for user received successfully")
}

func GetRecords(request *scheduleService.GetRecordsRequest, response *scheduleService.GetRecordsResponse) {
	records, err := postgres.GetRecords()
	if err != nil {
		logger.CreateLog("error", "get records: %v", err)
		response.Status = &scheduleService.Status{
			Successfully: false,
			Message:      "error getting records",
		}
		return
	}

	for records.Next() {
		var trainerId sql.NullInt64
		var dateStart, dateEnd sql.NullTime

		record := &scheduleService.Record{}
		err = records.Scan(&record.Id, &record.ClientId, &trainerId, &dateStart, &dateEnd)
		if err != nil {
			logger.CreateLog("error", "scan records: %v", err)
			response.Status = &scheduleService.Status{
				Successfully: false,
				Message:      "error scanning records",
			}
			return
		}

		if trainerId.Valid {
			record.TrainerId = trainerId.Int64
		}

		record.DateStart = timestamppb.New(dateStart.Time)
		record.DateEnd = timestamppb.New(dateEnd.Time)

		response.Records = append(response.Records, record)

	}
	response.Status = &scheduleService.Status{
		Successfully: true,
		Message:      "records received successfully",
	}
	logger.CreateLog("info", "records received successfully")
}

func GetRecordsByTrainerForDay(request *scheduleService.GetRecordsByTrainerForDayRequest, response *scheduleService.GetRecordsByTrainerForDayResponse) {
	date := utils.ConvertTimestampToString(request.Day)

	records, err := postgres.GetRecordsByTrainerForDay(request.TrainerId, date)
	if err != nil {
		logger.CreateLog("error", "get records by trainer for day: %v", err)
		response.Status = &scheduleService.Status{
			Successfully: false,
			Message:      "error getting records by trainer for day",
		}
		return
	}

	response.Time = make([]int64, 0)
	for records.Next() {
		var trainerId sql.NullInt64
		var dateStart, dateEnd sql.NullTime

		record := &scheduleService.Record{}
		err = records.Scan(&record.Id, &record.ClientId, &trainerId, &dateStart, &dateEnd)
		if err != nil {
			logger.CreateLog("error", "scan records by trainer for day: %v", err)
			response.Status = &scheduleService.Status{
				Successfully: false,
				Message:      "error scanning records by trainer for day",
			}
			return
		}

		response.Time = append(response.Time, int64(dateStart.Time.Hour()))
	}
	response.Status = &scheduleService.Status{
		Successfully: true,
		Message:      "records by trainer for day received successfully",
	}
	logger.CreateLog("info", "records by trainer for day received successfully")
}

func AddRecord(request *scheduleService.AddRecordRequest, response *scheduleService.AddRecordResponse) {
	var idTrainer sql.NullInt64
	idTrainer.Int64 = request.Record.TrainerId
	if idTrainer.Int64 != 0 {
		idTrainer.Valid = true
	}

	fmt.Println("idTrainer: ", idTrainer)

	_, err := postgres.AddRecord(request.Record.ClientId, idTrainer, request.Record.DateStart, request.Record.DateEnd)
	if err != nil {
		logger.CreateLog("error", "add record: %v", err)
		response.Status = &scheduleService.Status{
			Successfully: false,
			Message:      "error adding record",
		}
		return
	}

	logger.CreateLog("info", "record added successfully")
	response.Status = &scheduleService.Status{
		Successfully: true,
		Message:      "record added successfully",
	}
}

func DeleteRecordById(request *scheduleService.DeleteRecordByIdRequest, response *scheduleService.DeleteRecordByIdResponse) {
	_, err := postgres.DeleteRecordById(request.Id)
	if err != nil {
		logger.CreateLog("error", "delete record by id: %v", err)
		response.Status = &scheduleService.Status{
			Successfully: false,
			Message:      "error deleting record by id",
		}
		return
	}

	logger.CreateLog("info", "record deleted successfully")
	response.Status = &scheduleService.Status{
		Successfully: true,
		Message:      "record deleted successfully",
	}
}
