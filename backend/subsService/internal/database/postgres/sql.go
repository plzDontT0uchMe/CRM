package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Registration(idClient int64, idSubscription int64) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "INSERT INTO active (id_client, id_subscription) VALUES ($1, $2)", idClient, idSubscription)
}

func GetSubscriptions() (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT subscription.id, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities FROM subscription LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription GROUP BY subscription.id, subscription.name, subscription.price, subscription.description ORDER BY subscription.price")
}
func GetSubscriptionByAccountId(id int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT subscription.id, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities, active.id_trainer, active.date_expires FROM subscription LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription JOIN active ON subscription.id = active.id_subscription WHERE active.id_client = $1 GROUP BY subscription.id, subscription.name, subscription.price, subscription.description, active.id_trainer, active.date_expires", id)
}

func GetSubscriptionById(id int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT subscription.id, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities FROM subscription LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription WHERE subscription.id = $1 GROUP BY subscription.id, subscription.name, subscription.price, subscription.description", id)
}

func ChangeSubscription(idClient int64, idSubscription int64, idTrainer sql.NullInt64, dateExpires sql.NullTime) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE active SET id_subscription = $1, id_trainer = $2, date_expires = $3 WHERE id_client = $4", idSubscription, idTrainer, dateExpires, idClient)
}

func CreateApplication(idClient int64, idSubscription int64, idTrainer sql.NullInt64, duration sql.NullInt64) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "INSERT INTO applications (id_client, id_subcription, id_trainer, duration) VALUES ($1, $2, $3, $4) ON CONFLICT (id_client) DO UPDATE SET id_subcription = excluded.id_subcription, id_trainer = excluded.id_trainer, duration = excluded.duration", idClient, idSubscription, idTrainer, duration)
}

func GetApplications() (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT applications.id, applications.id_client, applications.id_trainer, applications.id_subcription, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities, applications.duration FROM applications LEFT JOIN subscription ON applications.id_subcription = subscription.id LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription GROUP BY applications.id, subscription.id, subscription.name, subscription.price, subscription.description ORDER BY applications.id")
}

func GetApplicationById(id int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT applications.id, applications.id_client, applications.id_trainer, applications.id_subcription, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities, applications.duration FROM applications LEFT JOIN subscription ON applications.id_subcription = subscription.id LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription WHERE applications.id = $1 GROUP BY applications.id, subscription.id, subscription.name, subscription.price, subscription.description", id)
}

func GetApplicationByAccountId(id int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT applications.id, applications.id_client, applications.id_trainer, applications.id_subcription, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities, applications.duration FROM applications LEFT JOIN subscription ON applications.id_subcription = subscription.id LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription WHERE applications.id_client = $1 GROUP BY applications.id, subscription.id, subscription.name, subscription.price, subscription.description", id)
}

func DeleteApplicationById(id int64) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "DELETE FROM applications WHERE id = $1", id)
}

func GetUsersByTrainerId(id int64) (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT id_client FROM active WHERE id_trainer = $1", id)
}
