package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetSubscriptions() (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT subscription.id, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities FROM subscription LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription GROUP BY subscription.id, subscription.name, subscription.price, subscription.description ORDER BY subscription.price")
}
func GetSubscriptionByAccountId(id int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT subscription.id, subscription.name, subscription.price, subscription.description, STRING_AGG(possibilities.possibility, ',') AS possibilities, active.id_trainer, active.date_expires FROM subscription LEFT JOIN possibilities ON subscription.id = possibilities.id_subscription JOIN active ON subscription.id = active.id_subscription WHERE active.id_client = $1 GROUP BY subscription.id, subscription.name, subscription.price, subscription.description, active.id_trainer, active.date_expires ORDER BY subscription.price", id)
}

func ChangeSubscription(idClient int64, idSubscription int64, idTrainer sql.NullInt64, dateExpires sql.NullTime) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE active SET id_subscription = $1, id_trainer = $2, date_expires = $3 WHERE id_client = $4", idSubscription, idTrainer, dateExpires, idClient)
}
