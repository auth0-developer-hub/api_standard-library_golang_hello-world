package mock

import (
	"context"
	"hello-golang-api/entities"

	"github.com/snowzach/queryp"
)

type Client struct {
	//MessageGetByIDfn mock method to use during testing
	MessageGetByIDfn func(ctx context.Context, id string) (*entities.Message, error)
	//MessageSavefn mock method to use during testing
	MessageSavefn func(ctx context.Context, user *entities.Message) error

	MessageUpdatefn func(ctx context.Context, user *entities.Message) error
	//MessageDeleteByIDfn mock method to use during testing
	MessageDeleteByIDfn func(ctx context.Context, id string) error
	//MessagesListfn mock method to use during testing
	MessagesListfn func(ctx context.Context, qp *queryp.QueryParameters) ([]*entities.Message, int64, error)
}

func (c *Client) Close() error {
	return nil
}
func (c *Client) MessageGetByID(ctx context.Context, id string) (*entities.Message, error) {
	return c.MessageGetByIDfn(ctx, id)
}
func (c *Client) MessageSave(ctx context.Context, user *entities.Message) error {
	return c.MessageSavefn(ctx, user)
}
func (c *Client) MessageDeleteByID(ctx context.Context, id string) error {
	return c.MessageDeleteByIDfn(ctx, id)
}
func (c *Client) MessagesList(ctx context.Context, qp *queryp.QueryParameters) ([]*entities.Message, int64, error) {
	return c.MessagesListfn(ctx, qp)
}
func (c *Client) MessageUpdate(ctx context.Context, user *entities.Message) error {
	return c.MessageUpdatefn(ctx, user)
}
