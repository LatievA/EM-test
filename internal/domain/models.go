package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price"        db:"price"`
	UserID      uuid.UUID `json:"user_id"      db:"user_id"`

	// Internally stored as time.Time (DATE in PostgreSQL) // may be deleted
	StartDate time.Time  `json:"-"             db:"start_date"`
	EndDate   *time.Time `json:"-"             db:"end_date"`

	// For JSON input/output (MM-YYYY format)
	StartDateString string  `json:"start_date"`
	EndDateString   *string `json:"end_date,omitempty"`
}

// ParseDates parses MM-YYYY strings into time.Time
func (s *Subscription) ParseDates() error {
	var err error
	s.StartDate, err = time.Parse("01-2006", s.StartDateString)
	if err != nil {
		return err
	}
	if s.EndDateString != nil {
		endDate, err := time.Parse("01-2006", *s.EndDateString)
		if err != nil {
			return err
		}
		s.EndDate = &endDate
	}
	return nil
}

// FormatDates formats time.Time into MM-YYYY strings
func (s *Subscription) FormatDates() {
	s.StartDateString = s.StartDate.Format("01-2006")
	if s.EndDate != nil {
		endDateStr := s.EndDate.Format("01-2006")
		s.EndDateString = &endDateStr
	}
}
