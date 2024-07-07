package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetRecordsForUser(idClient int64) (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT * FROM records WHERE id_client = $1 OR id_trainer = $1", idClient)
}

func GetRecords() (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT * FROM records")
}

func GetRecordsByTrainerForDay(idTrainer int64, date string) (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT * FROM records WHERE id_trainer = $1 AND DATE(date_start) = $2", idTrainer, date)
}

func AddRecord(idClient int64, idTrainer sql.NullInt64, dateStart *timestamppb.Timestamp, dateEnd *timestamppb.Timestamp) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "INSERT INTO records (id_client, id_trainer, date_start, date_end) VALUES ($1, $2, $3, $4)", idClient, idTrainer, dateStart.AsTime(), dateEnd.AsTime())
}

func DeleteRecordById(idRecord int64) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "DELETE FROM records WHERE id = $1", idRecord)
}
