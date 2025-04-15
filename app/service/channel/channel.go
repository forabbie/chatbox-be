package service

import (
	"chatbox/pkg/database"
	"context"

	mchannel "chatbox/app/model/channel"
)

func Insert(ctx context.Context, name string, createdBy int64, userIDs []int64) (*mchannel.Channel, error) {
	tx, err := database.PostgresMain.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Insert into channels
	var channelID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO channels (name, created_by)
		VALUES ($1, $2)
		RETURNING id
	`, name, createdBy).Scan(&channelID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert creator as admin
	_, err = tx.ExecContext(ctx, `
		INSERT INTO channel_members (channel_id, user_id, role)
		VALUES ($1, $2, 'admin')
	`, channelID, createdBy)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert other members as default 'member'
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO channel_members (channel_id, user_id, role)
		VALUES ($1, $2, 'member')
		ON CONFLICT (channel_id, user_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	for _, userID := range userIDs {
		if userID == createdBy {
			continue // skip if already added as admin
		}
		if _, err := stmt.ExecContext(ctx, channelID, userID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &mchannel.Channel{
		ID:      channelID,
		Name:    name,
		UserIDs: append([]int64{createdBy}, userIDs...), // full member list
	}, nil
}
