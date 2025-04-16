package service

import (
	"chatbox/pkg/database"
	"context"

	mchannel "chatbox/app/model/channel"
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

// func GetByUserID(ctx context.Context, userID int64) ([]*mchannel.ChannelParam, error) {
// 	query := ` SELECT c.id, c.name, c.created_by
// 		FROM channels c
// 		INNER JOIN channel_members cu ON cu.channel_id = c.id
// 		WHERE cu.user_id = $1`

// 	// row := database.PostgresMain.DB.QueryRowContext(ctx, query, userID)
// 	// if err := row.Err(); err != nil {
// 	// 	return nil, err
// 	// }

// 	// channels := new(mchannel.ChannelParam)

// 	// if err := row.Scan(&channels.ID, &channels.Name, &channels.CreatedBy); err != nil {
// 	// 	return nil, err
// 	// }

// 	// defer rows.Close()

// 	// var channels []ChannelParam

// 	// for rows.Next() {
// 	// 	var ch ChannelParam
// 	// 	if err := rows.Scan(&ch.ID, &ch.Name, &ch.CreatedBy); err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	channels = append(channels, ch)
// 	// }

// 	return []*mchannel.ChannelParam{channels}, nil
// }

func GetByUserID(ctx context.Context, userID int64) ([]*mchannel.ChannelParam, error) {
	query := ` SELECT c.id, c.name, c.created_by
		FROM channels c
		INNER JOIN channel_members cu ON cu.channel_id = c.id
		WHERE cu.user_id = $1`

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*mchannel.ChannelParam
	for rows.Next() {
		var ch mchannel.ChannelParam
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.CreatedBy); err != nil {
			return nil, err
		}
		channels = append(channels, &ch)
	}

	return channels, nil
}
