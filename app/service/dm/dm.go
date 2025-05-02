package service

import (
	"chatbox/pkg/database"
	"context"

	mdm "chatbox/app/model/dm"
)

func GetDMListByUserID(ctx context.Context, userID int64) ([]*mdm.DMListItem, error) {
	query := `
		WITH latest_messages AS (
			SELECT DISTINCT ON (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id))
				id,
				sender_id,
				receiver_id,
				message,
				sent_at
			FROM direct_messages
			WHERE sender_id = $1 OR receiver_id = $1
			ORDER BY LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id), sent_at DESC
		)
		SELECT 
			id,
			sender_id,
			receiver_id,
			message,
			sent_at
		FROM latest_messages
		ORDER BY sent_at DESC
	`

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*mdm.DMListItem
	for rows.Next() {
		var item mdm.DMListItem
		if err := rows.Scan(&item.ID, &item.SenderID, &item.ReceiverID, &item.Message, &item.SentAt); err != nil {
			return nil, err
		}
		results = append(results, &item)
	}

	return results, nil
}
