package service

import (
	"context"
	"em-test/internal/adapter/db"
	"em-test/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo *db.SubscriptionRepo
}

func NewSubscriptionService(repo *db.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	sub *domain.Subscription,
) error {
	if sub.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	if sub.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	if sub.StartDateString == "" {
		return fmt.Errorf("start_date is required")
	}

	if err := sub.ParseDates(); err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) GetSubscription(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
) (*domain.Subscription, error) {
	return s.repo.Read(ctx, userID, serviceName)
}

func (s *SubscriptionService) UpdateSubscription(
	ctx context.Context,
	sub *domain.Subscription,
) error {
	if sub.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	if sub.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	if sub.StartDateString == "" {
		return fmt.Errorf("start_date is required")
	}

	if err := sub.ParseDates(); err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) DeleteSubscription(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
) error {
	return s.repo.Delete(ctx, userID, serviceName)
}

// List subscriptions for a user
func (s *SubscriptionService) ListSubscriptions(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Subscription, error) {
	return s.repo.List(ctx, userID)
}

func (s *SubscriptionService) SumSubscriptionsPrice(
	ctx context.Context,
	userID uuid.UUID,
	serviceName *string,
	startDateString string,
	endDateString string,
) (int, error) {
	startDate, err := time.Parse("01-2006", startDateString)
	if err != nil {
		return 0, fmt.Errorf("invalid start_date format: %v", err)
	}
	endDate, err := time.Parse("01-2006", endDateString)
	if err != nil {
		return 0, fmt.Errorf("invalid end_date format: %v", err)
	}

	// Use first day of the next month to include the entire specified month
	endDate = endDate.AddDate(0, 1, 0)

	return s.repo.SumSubscriptionsPrice(ctx, userID, serviceName, startDate, endDate)
}
