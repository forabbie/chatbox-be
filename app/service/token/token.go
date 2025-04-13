package service

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	mtoken "chatbox/app/model/token"

// 	"chatbox/pkg/database"
// )

// func Insert(ctx context.Context, token *mtoken.Token) (int64, error) {
// 	query := "INSERT "

// 	query += "INTO token(id, user_id, value, type, expiration) "

// 	query += "VALUES(?, ?, ?, ?, ?)"

// 	res, err := database.MySQLd.DB.ExecContext(
// 		ctx,
// 		query,
// 		token.Id,
// 		token.UserId,
// 		token.Value,
// 		token.Type,
// 		token.Expiration,
// 	)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return res.LastInsertId()
// }

// func Fetch(ctx context.Context, filter map[string][]string, args []interface{}, order, sort string, limit, offset int) ([]mtoken.Token, error) {
// 	query := "SELECT * "

// 	query += "FROM token "

// 	if len(filter) > 0 {
// 		query += "WHERE "

// 		field := []string{}

// 		if _, ok := filter["or"]; ok {
// 			field = append(field, (strings.Join(filter["or"], " OR ")))
// 		}

// 		if _, ok := filter["and"]; ok {
// 			field = append(field, (strings.Join(filter["and"], " AND ")))
// 		}

// 		query += strings.Join(field, " AND ") + " "
// 	}

// 	query += fmt.Sprintf("ORDER BY %s %s ", order, sort)

// 	query += "LIMIT ? OFFSET ?"

// 	args = append(args, limit, offset)

// 	rows, err := database.MySQLd.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tokens := []mtoken.Token{}

// 	for rows.Next() {
// 		token := mtoken.Token{}

// 		if err := rows.Scan(
// 			&token.Id,
// 			&token.UserId,
// 			&token.Value,
// 			&token.Type,
// 			&token.Expiration,
// 			&token.Created,
// 			&token.LastModified,
// 		); err != nil {
// 			return nil, err
// 		}

// 		tokens = append(tokens, token)
// 	}

// 	if err := rows.Close(); err != nil {
// 		return nil, err
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return tokens, nil
// }
