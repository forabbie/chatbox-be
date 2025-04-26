package service

import (
	"chatbox/pkg/database"
	"context"
	"fmt"
	"log"
	"strings"

	mmsg "chatbox/app/model/message"

	"chatbox/pkg/util"
)

func Insert(ctx context.Context, msg *mmsg.Message) (*mmsg.Message, error) {
	var query string

	if msg.ReceiverClass == "user" {
		query = `
			INSERT INTO direct_messages (sender_id, receiver_id, message)
			VALUES ($1, $2, $3)
			RETURNING id, sent_at
		`
	} else if msg.ReceiverClass == "channel" {
		query = `
			INSERT INTO channel_messages (sender_id, channel_id, message)
			VALUES ($1, $2, $3)
			RETURNING id, sent_at
		`
	} else {
		return nil, fmt.Errorf("invalid receiver_class: %s", msg.ReceiverClass)
	}

	// Execute the query
	var err error
	if msg.ReceiverClass == "user" {
		err = database.PostgresMain.DB.QueryRowContext(
			ctx,
			query,
			msg.Sender.ID,
			msg.ReceiverID,
			msg.Message,
		).Scan(&msg.ID, &msg.SentAt)
	} else if msg.ReceiverClass == "channel" {
		// Channel message
		err = database.PostgresMain.DB.QueryRowContext(
			ctx,
			query,
			msg.Sender.ID,
			msg.ReceiverID,
			msg.Message,
		).Scan(&msg.ID, &msg.SentAt)
	}

	// Check for errors
	if err != nil {
		log.Printf("Database error: %v, Query: %v", err, query)
		return nil, err
	}

	return msg, nil
}

func Fetch(ctx context.Context, filter map[string][]string, args []interface{}, order, sort string, limit, offset int) ([]mmsg.DirectMessage, error) {
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
				dm.id AS message_id,
				dm.message AS message,
				dm.sent_at,
				dm.is_edited,
				dm.edited_at,
				dm.deleted_at,
				sender.id AS sender_id,
				sender.username AS sender_username,
				sender.firstname AS sender_firstname,
				sender.lastname AS sender_lastname,
				receiver.id AS receiver_id,
				receiver.username AS receiver_username,
				receiver.firstname AS receiver_firstname,
				receiver.lastname AS receiver_lastname
		FROM direct_messages dm
		JOIN users sender ON sender.id = dm.sender_id
		JOIN users receiver ON receiver.id = dm.receiver_id
		WHERE dm.receiver_class = 'user' AND (
				(dm.sender_id = ? AND dm.receiver_id = ?) OR 
				(dm.receiver_id = ? AND dm.sender_id = ?)
		)
		UNION
		SELECT
				chm.id AS message_id,
				chm.message AS message,
				chm.sent_at,
				chm.is_edited,
				chm.edited_at,
				chm.deleted_at,
				sender.id AS sender_id,
				sender.username AS sender_username,
				sender.firstname AS sender_firstname,
				sender.lastname AS sender_lastname,
				null AS receiver_id,
				null AS receiver_username,
				null AS receiver_firstname,
				null AS receiver_lastname
		FROM channel_messages chm
		JOIN users sender ON sender.id = chm.sender_id
		WHERE chm.channel_id = ?
		ORDER BY sent_at DESC
		LIMIT ? OFFSET ?
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

	var messages []mmsg.DirectMessage
	for rows.Next() {
		var msg mmsg.DirectMessage
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
