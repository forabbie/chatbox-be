package service

import (
	"context"
	"fmt"
	"strings"

	muser "chatbox/app/model/user"

	"chatbox/pkg/database"
)

func Insert(ctx context.Context, user *muser.User) (int64, error) {
	query := `
		INSERT INTO "users" (firstname, lastname, username, emailaddress, hashed_password, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := database.PostgresMain.DB.QueryRowContext(
		ctx,
		query,
		user.Firstname,
		user.Lastname,
		user.Username,
		user.EmailAddress,
		user.Password,
		user.IsActive,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func Count(ctx context.Context, filter map[string][]string, args []interface{}) (int64, error) {
	query := `SELECT COUNT(*) FROM users`
	conditions := []string{}
	paramIndex := 1

	if orFilters, ok := filter["or"]; ok && len(orFilters) > 0 {
		orConditions := []string{}
		for range orFilters {
			orConditions = append(orConditions, fmt.Sprintf("emailaddress = $%d", paramIndex))
			paramIndex++
		}
		conditions = append(conditions, "("+strings.Join(orConditions, " OR ")+")")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := database.PostgresMain.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
