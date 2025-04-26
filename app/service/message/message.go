package service

import (
	"chatbox/pkg/database"
	"context"
	"fmt"
	"strings"

	mdmsg "chatbox/app/model/message"

	"chatbox/pkg/util"
)

func Insert(ctx context.Context, dmsg *mdmsg.DirectMessage) (*mdmsg.DirectMessage, error) {
	query := `INSERT INTO direct_messages (sender_id, receiver_id, message)
	          VALUES ($1, $2, $3) RETURNING id, sent_at`

	err := database.PostgresMain.DB.QueryRowContext(ctx, query, dmsg.Sender.ID, dmsg.Receiver.ID, dmsg.Message).
		Scan(&dmsg.ID, &dmsg.SentAt)

	if err != nil {
		return nil, err
	}

	return dmsg, nil
}

func Fetch(ctx context.Context, filter map[string][]string, args []interface{}, order, sort string, limit, offset int) ([]mdmsg.DirectMessage, error) {
	allowedOrderFields := map[string]bool{
		"dm.sent_at":    true,
		"dm.edited_at":  true,
		"dm.deleted_at": true,
	}
	if !allowedOrderFields[order] {
		order = "dm.sent_at"
	}

	sort = strings.ToUpper(sort)
	if sort != "DESC" {
		sort = "ASC"
	}

	query := `
		SELECT
			dm.id,
			dm.message,
			dm.sent_at,
			dm.is_edited,
			dm.edited_at,
			dm.deleted_at,

			sender.id, sender.username, sender.firstname, sender.lastname,
			receiver.id, receiver.username, receiver.firstname, receiver.lastname

		FROM direct_messages dm
		JOIN users sender ON sender.id = dm.sender_id
		JOIN users receiver ON receiver.id = dm.receiver_id
	`

	var conditions []string

	if ors, ok := filter["or"]; ok && len(ors) > 0 {
		conditions = append(conditions, "("+strings.Join(ors, " OR ")+")")
	}

	if ands, ok := filter["and"]; ok && len(ands) > 0 {
		conditions = append(conditions, "("+strings.Join(ands, " AND ")+")")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	query, err := util.ReplacePlaceholders(query, len(args))
	if err != nil {
		return nil, err
	}

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []mdmsg.DirectMessage
	for rows.Next() {
		var msg mdmsg.DirectMessage
		if err := rows.Scan(
			&msg.ID,
			&msg.Message,
			&msg.SentAt,
			&msg.IsEdited,
			&msg.EditedAt,
			&msg.DeletedAt,
			&msg.Sender.ID,
			&msg.Sender.Username,
			&msg.Sender.Firstname,
			&msg.Sender.Lastname,
			&msg.Receiver.ID,
			&msg.Receiver.Username,
			&msg.Receiver.Firstname,
			&msg.Receiver.Lastname,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
