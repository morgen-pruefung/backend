package newsletterstore

import (
	"backend/internal/analytics"
	"backend/internal/newsletter"
	"context"
	"fmt"
	"regexp"
)

type DB interface {
	Subscribe(ctx context.Context, e newsletter.Entry) error
	Unsubscribe(ctx context.Context, e newsletter.Entry) error
	GetSubscribers(ctx context.Context) ([]newsletter.Entry, error)
}

const EmailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

type Store struct {
	db DB
}

func NewStore(db DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Subscribe(ctx context.Context, e newsletter.Entry) error {
	if e.Email == "" {
		return newsletter.ErrEmptyEmail
	}

	if matched, _ := regexp.MatchString(EmailRegex, e.Email); !matched {
		return newsletter.ErrorInvalidEmail
	}

	if err := s.db.Subscribe(ctx, e); err != nil {
		return fmt.Errorf("could not subscribe to newsletter: %w", err)
	}

	analytics.SendEvent(analytics.Event{
		Name:       "NewsletterSubscribed",
		Properties: map[string]interface{}{},
	})

	return nil
}

func (s *Store) Unsubscribe(ctx context.Context, e newsletter.Entry) error {
	err := s.db.Unsubscribe(ctx, e)
	if err != nil {
		return fmt.Errorf("could not unsubscribe from newsletter: %w", err)
	}

	analytics.SendEvent(analytics.Event{
		Name:       "NewsletterUnsubscribed",
		Properties: map[string]interface{}{},
	})

	return nil
}

func (s *Store) GetSubscribers(ctx context.Context) ([]newsletter.Entry, error) {
	return s.db.GetSubscribers(ctx)
}
