package service

import (
	"chatbox/pkg/database"
	"context"

	mchannel "chatbox/app/model/channel"

	"github.com/lib/pq"
)

func Insert(ctx context.Context, name string, createdBy int64, userIDs []int64) (*mchannel.ChannelParam, error) {
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

	return &mchannel.ChannelParam{
		ID:        channelID,
		Name:      name,
		CreatedBy: createdBy,
		UserIDs:   append([]int64{createdBy}, userIDs...),
	}, nil
}

func GetByUserID(ctx context.Context, userID int64) ([]*mchannel.ChannelParam, error) {
	query := `
		SELECT 
			c.id,
			c.name,
			c.created_by,
			ARRAY_AGG(cm_all.user_id) AS user_ids
		FROM channels c
		-- join only to filter channels where user is a member
		JOIN channel_members cm_filter ON cm_filter.channel_id = c.id AND cm_filter.user_id = $1
		-- join again to get all members
		JOIN channel_members cm_all ON cm_all.channel_id = c.id
		GROUP BY c.id, c.name, c.created_by
	`

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*mchannel.ChannelParam
	for rows.Next() {
		var ch mchannel.ChannelParam
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.CreatedBy, pq.Array(&ch.UserIDs)); err != nil {
			return nil, err
		}
		channels = append(channels, &ch)
	}

	return channels, nil
}
