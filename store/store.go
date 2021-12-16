package store

import (
	"context"
	"hello-golang-api/entities"

	"github.com/snowzach/queryp"
)

// MessageStore is the persistent store of messages
type MessageStore interface {
	MessageGetByID(ctx context.Context, id string) (*entities.Message, error)
	MessageSave(ctx context.Context, user *entities.Message) error
	MessageUpdate(ctx context.Context, user *entities.Message) error
	MessageDeleteByID(ctx context.Context, id string) error
	MessagesList(ctx context.Context, qp *queryp.QueryParameters) ([]*entities.Message, int64, error)
}
