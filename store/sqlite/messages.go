package sqlite

import (
	"context"
	"fmt"
	"hello-golang-api/entities"
	"strconv"
	"strings"

	"github.com/snowzach/queryp"
	"github.com/snowzach/queryp/qppg"
)

func (c *Client) MessageGetByID(ctx context.Context, id string) (*entities.Message, error) {
	msg := new(entities.Message)
	row := c.db.QueryRow(`SELECT id, text, date FROM messages WHERE id=?`, id)

	err := row.Scan(&msg.Id, &msg.Text, &msg.Date)
	if err != nil {
		return nil, fmt.Errorf("no record found for %q", id)
	}
	return msg, nil
}

func (c *Client) MessageSave(ctx context.Context, msg *entities.Message) error {
	stmt, err := c.db.Prepare("insert into messages(id,text,date) values(?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(msg.Id, msg.Text, msg.Date)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	c.logger.Debugf("updated %d records", affect)
	return nil
}

func (c *Client) MessageUpdate(ctx context.Context, user *entities.Message) error {
	stmt, err := c.db.Prepare("update messages set text=? where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(user.Text, user.Id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	c.logger.Debugf("updated %d records", affect)
	return nil
}

func (c *Client) MessageDeleteByID(ctx context.Context, id string) error {

	// delete
	stmt, err := c.db.Prepare("DELETE from messages where id=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	c.logger.Debugf("Deleted %d records", affect)
	return nil
}

func (c *Client) MessagesList(ctx context.Context, qp *queryp.QueryParameters) ([]*entities.Message, int64, error) {
	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := queryp.FilterFieldTypes{
		"messages.id":   queryp.FilterTypeSimple,
		"messages.text": queryp.FilterTypeString,
		"messages.date": queryp.FilterTypeString,
	}

	sortFields := queryp.SortFields{
		"messages.id":   "",
		"messages.text": "",
		"messages.date": "",
	}
	// Default sort
	if len(qp.Sort) == 0 {
		qp.Sort.Append("messages.id", false)
	}

	if len(qp.Filter) > 0 {
		queryClause.WriteString(" WHERE ")
	}

	if err := qppg.FilterQuery(filterFields, qp.Filter, &queryClause, &queryParams); err != nil {
		return nil, 0, err
	}
	var count int64
	if row := c.db.QueryRow(`SELECT COUNT(*) AS count FROM messages `+queryClause.String(), queryParams...); row != nil {
		err := row.Scan(&count)
		if err != nil {
			return nil, 0, err
		}
	}
	if err := qppg.SortQuery(sortFields, qp.Sort, &queryClause, &queryParams); err != nil {
		return nil, 0, err
	}
	if qp.Limit > 0 {
		queryClause.WriteString(" LIMIT " + strconv.FormatInt(qp.Limit, 10))
	}
	if qp.Offset > 0 {
		queryClause.WriteString(" OFFSET " + strconv.FormatInt(qp.Offset, 10))
	}

	var records = make([]*entities.Message, 0)
	rows, err := c.db.Query(`SELECT id, text, date FROM messages`+queryClause.String(), queryParams...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := &entities.Message{}
		err = rows.Scan(&msg.Id, &msg.Text, &msg.Date)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, msg)
	}

	return records, count, nil
}
