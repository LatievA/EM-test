package db

import (
	"context"
	"time"

	"em-test/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepo(ctx context.Context, dsn string) (*SubscriptionRepo, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return &SubscriptionRepo{pool: pool}, nil
}

func (r *SubscriptionRepo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	if err := sub.ParseDates(); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)
	return err
}

func (r *SubscriptionRepo) Read(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND service_name = $2
	`, userID, serviceName).Scan(
		&sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate,
	)
	if err != nil {
		return nil, err
	}
	sub.FormatDates()
	return &sub, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	if err := sub.ParseDates(); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE subscriptions
		SET price = $3, start_date = $4, end_date = $5
		WHERE user_id = $1 AND service_name = $2
	`, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)
	return err
}

func (r *SubscriptionRepo) Delete(ctx context.Context, userID uuid.UUID, serviceName string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM subscriptions
		WHERE user_id = $1 AND service_name = $2
	`, userID, serviceName)
	return err
}

func (r *SubscriptionRepo) List(ctx context.Context, userID uuid.UUID) ([]domain.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		if err := rows.Scan(
			&sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate,
		); err != nil {
			return nil, err
		}
		sub.FormatDates()
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepo) SumSubscriptionsPrice(
    ctx context.Context,
    userID uuid.UUID,
    serviceName *string,
    startDate time.Time,
    endDate time.Time,
) (int, error) {
    query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE user_id = $1
        AND start_date <= $3
        AND (end_date IS NULL OR end_date >= $2)
    `
    args := []interface{}{userID, startDate, endDate}

    if serviceName != nil {
        query += " AND service_name = $4"
        args = append(args, *serviceName)
    }

    var totalPrice int
    err := r.pool.QueryRow(ctx, query, args...).Scan(&totalPrice)
    return totalPrice, err
}