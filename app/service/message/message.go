package service

import (
	"chatbox/pkg/database"
	"context"

	mdmsg "chatbox/app/model/message"
)

func Insert(ctx context.Context, dmsg *mdmsg.DirectMessage) (*mdmsg.DirectMessage, error) {
	query := `INSERT INTO direct_messages (sender_id, receiver_id, message)
	          VALUES ($1, $2, $3) RETURNING id, sent_at`

	err := database.PostgresMain.DB.QueryRowContext(ctx, query, dmsg.SenderID, dmsg.ReceiverID, dmsg.Message).
		Scan(&dmsg.ID, &dmsg.SentAt)

	if err != nil {
		return nil, err
	}

	return dmsg, nil
}
