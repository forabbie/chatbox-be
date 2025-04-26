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

func FetchDirectMessages(ctx context.Context, userID, receiverID int64, filter map[string][]string, args []interface{}, order, sort string, limit, offset int) ([]mmsg.Message, error) {
	query := `
			SELECT
					dm.id, dm.message, dm.sent_at, dm.is_edited, dm.edited_at, dm.deleted_at,
					sender.id, sender.username, sender.firstname, sender.lastname,
					receiver.id, receiver.username, receiver.firstname, receiver.lastname
			FROM direct_messages dm
			JOIN users sender ON sender.id = dm.sender_id
			JOIN users receiver ON receiver.id = dm.receiver_id
	`

	// Append filters
	if q := strings.Join(filter["and"], " AND "); q != "" {
		query += " WHERE " + q
	}

	if q := strings.Join(filter["or"], " OR "); q != "" {
		if !strings.Contains(query, "WHERE") {
			query += " WHERE "
		} else {
			query += " AND "
		}
		query += "(" + q + ")"
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", order, sort)
	args = append(args, limit, offset)

	query, _ = util.ReplacePlaceholders(query, len(args))

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []mmsg.Message
	for rows.Next() {
		var msg mmsg.Message
		msg.Receiver = &mmsg.User{}

		err := rows.Scan(
			&msg.ID, &msg.Message, &msg.SentAt, &msg.IsEdited, &msg.EditedAt, &msg.DeletedAt,
			&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.Firstname, &msg.Sender.Lastname,
			&msg.Receiver.ID, &msg.Receiver.Username, &msg.Receiver.Firstname, &msg.Receiver.Lastname,
		)
		if err != nil {
			return nil, err
		}

		msg.ReceiverClass = "user"
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func FetchChannelMessages(ctx context.Context, channelID int64, filter map[string][]string, args []interface{}, order, sort string, limit, offset int) ([]mmsg.Message, error) {
	query := `
		SELECT
			chm.id, chm.message, chm.sent_at, chm.is_edited, chm.edited_at, chm.deleted_at,
			sender.id, sender.username, sender.firstname, sender.lastname,
			NULL, NULL, NULL, NULL
		FROM channel_messages chm
		JOIN users sender ON sender.id = chm.sender_id
		WHERE chm.channel_id = ?
	`

	args = append(args, channelID)

	if q := strings.Join(filter["and"], " AND "); q != "" {
		query += " AND " + q
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", order, sort)
	args = append(args, limit, offset)

	query, _ = util.ReplacePlaceholders(query, len(args))

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []mmsg.Message
	for rows.Next() {
		var msg mmsg.Message

		// Temporary null fields for receiver
		var recvID *int64
		var recvUsername, recvFirstname, recvLastname *string

		err := rows.Scan(
			&msg.ID, &msg.Message, &msg.SentAt, &msg.IsEdited, &msg.EditedAt, &msg.DeletedAt,
			&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.Firstname, &msg.Sender.Lastname,
			&recvID, &recvUsername, &recvFirstname, &recvLastname,
		)
		if err != nil {
			return nil, err
		}

		// Receiver is nil in channel messages (by design)
		msg.Receiver = nil
		msg.ReceiverID = &channelID
		msg.ReceiverClass = "channel"

		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
