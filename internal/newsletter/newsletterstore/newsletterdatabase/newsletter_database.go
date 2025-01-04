package newsletterdatabase

import (
	"backend/internal/newsletter"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	collection *mongo.Collection
}

func NewDB(collection *mongo.Collection) *DB {
	return &DB{
		collection: collection,
	}
}

func (d *DB) Subscribe(ctx context.Context, e newsletter.Entry) error {
	_, err := d.collection.InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("could not subscribe to newsletter: %w", err)
	}

	return nil
}

func (d *DB) Unsubscribe(ctx context.Context, e newsletter.Entry) error {
	_, err := d.collection.DeleteOne(ctx, e)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil // no subscriber found
		}

		return fmt.Errorf("could not unsubscribe from newsletter: %w", err)
	}

	return nil
}

func (d *DB) GetSubscribers(ctx context.Context) ([]newsletter.Entry, error) {
	cursor, err := d.collection.Find(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get subscribers: %w", err)
	}
	defer cursor.Close(ctx)

	var subscribers []newsletter.Entry
	for cursor.Next(ctx) {
		var entry newsletter.Entry
		if err := cursor.Decode(&entry); err != nil {
			return nil, fmt.Errorf("could not decode subscriber: %w", err)
		}

		subscribers = append(subscribers, entry)
	}

	return subscribers, nil
}
