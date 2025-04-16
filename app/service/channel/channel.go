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

func GetDetailsByID(ctx context.Context, channelID int64) (*mchannel.GetChannelParam, error) {
	// Step 1: Get channel details
	query := `
		SELECT c.id, c.name, c.created_by
		FROM channels c
		WHERE c.id = $1
	`
	row := database.PostgresMain.DB.QueryRowContext(ctx, query, channelID)

	var ch mchannel.GetChannelParam
	if err := row.Scan(&ch.ID, &ch.Name, &ch.CreatedBy); err != nil {
		return nil, err
	}

	// Step 2: Get channel members
	memberQuery := `
		SELECT u.id, u.firstname, u.lastname, u.emailaddress
		FROM users u
		JOIN channel_members cm ON cm.user_id = u.id
		WHERE cm.channel_id = $1
	`

	rows, err := database.PostgresMain.DB.QueryContext(ctx, memberQuery, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []mchannel.UserSummary
	for rows.Next() {
		var u mchannel.UserSummary
		if err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Emailaddress); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	ch.Members = members

	// Step 3: Get creator's user details
	creatorQuery := `
		SELECT id, firstname, lastname, emailaddress
		FROM users
		WHERE id = $1
	`
	var creator mchannel.UserSummary
	if err := database.PostgresMain.DB.QueryRowContext(ctx, creatorQuery, ch.CreatedBy).Scan(
		&creator.ID, &creator.Firstname, &creator.Lastname, &creator.Emailaddress,
	); err != nil {
		return nil, err
	}
	ch.CreatedByUser = creator

	return &ch, nil
}

func AddMember(channelID int64, userID int64) error {
	_, err := database.PostgresMain.DB.Exec(`
		INSERT INTO channel_members (channel_id, user_id)
		VALUES ($1, $2)
	`, channelID, userID)

	return err
}

func IsMember(channelID int64, userID int64) (bool, error) {
	var exists bool
	err := database.PostgresMain.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM channel_members WHERE channel_id = $1 AND user_id = $2
		)
	`, channelID, userID).Scan(&exists)

	return exists, err
}

func IsChannelCreator(ctx context.Context, channelID, userID int64) (bool, error) {
	var createdBy int64
	err := database.PostgresMain.DB.QueryRowContext(ctx, `
		SELECT created_by FROM channels WHERE id = $1
	`, channelID).Scan(&createdBy)

	if err != nil {
		return false, err
	}
	return createdBy == userID, nil
}

func RemoveMember(ctx context.Context, channelID, userID int64) error {
	_, err := database.PostgresMain.DB.ExecContext(ctx, `
		DELETE FROM channel_members WHERE channel_id = $1 AND user_id = $2
	`, channelID, userID)
	return err
}

func Delete(ctx context.Context, channelID int64) error {
	_, err := database.PostgresMain.DB.ExecContext(ctx, `
		DELETE FROM channels WHERE id = $1
	`, channelID)
	return err
}
