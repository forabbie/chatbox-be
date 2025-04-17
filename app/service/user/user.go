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

func GetByID(ctx context.Context, id int64) (*muser.UserDetails, error) {
	query := "SELECT id, firstname, lastname, username, emailaddress, is_active "

	query += "FROM users "

	query += "WHERE id = $1"

	row := database.PostgresMain.DB.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	user := new(muser.UserDetails)

	if err := row.Scan(
		&user.Id,
		&user.Firstname,
		&user.Lastname,
		&user.Username,
		&user.EmailAddress,
		&user.IsActive,
	); err != nil {
		return nil, err
	}

	return user, nil
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

func Fetch(ctx context.Context, filter map[string][]string, args []interface{}, limit int) ([]muser.User, error) {
	query := "SELECT id, firstname, lastname, username, emailaddress, hashed_password, is_active FROM users"

	var conditions []string
	placeholderIndex := 1
	var finalArgs []interface{}

	// OR conditions
	if ors, ok := filter["or"]; ok && len(ors) > 0 {
		var orConditions []string
		for _, condition := range ors {
			orConditions = append(orConditions, fmt.Sprintf("%s $%d", strings.Split(condition, " ?")[0], placeholderIndex))
			placeholderIndex++
		}
		conditions = append(conditions, "("+strings.Join(orConditions, " OR ")+")")
		finalArgs = append(finalArgs, args[:len(ors)]...)
	}

	// AND conditions
	if ands, ok := filter["and"]; ok && len(ands) > 0 {
		var andConditions []string
		for _, condition := range ands {
			andConditions = append(andConditions, fmt.Sprintf("%s $%d", strings.Split(condition, " ?")[0], placeholderIndex))
			placeholderIndex++
		}
		conditions = append(conditions, "("+strings.Join(andConditions, " AND ")+")")
		finalArgs = append(finalArgs, args[len(finalArgs):]...)
	}

	// WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// LIMIT placeholder
	query += fmt.Sprintf(" LIMIT $%d", placeholderIndex)
	finalArgs = append(finalArgs, limit)

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []muser.User
	for rows.Next() {
		var user muser.User
		if err := rows.Scan(
			&user.Id,
			&user.Firstname,
			&user.Lastname,
			&user.Username,
			&user.EmailAddress,
			&user.Password,
			&user.IsActive,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetAll(ctx context.Context) ([]*muser.User, error) {
	query := `SELECT id, firstname, lastname, emailaddress FROM users`

	rows, err := database.PostgresMain.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*muser.User

	for rows.Next() {
		var user muser.User
		if err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.EmailAddress); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
