package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"mini-app-notifications/internal/domain"

	sl "mini-app-notifications/internal/logger"

	"github.com/Masterminds/squirrel"
)

type UserReposirory struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserReposirory {
	return &UserReposirory{
		db:     db,
		logger: logger,
	}
}

func (u *UserReposirory) GetUsers() ([]domain.User, error) {
	const op = "repository.user.GetUsers"

	sql, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "tg_user_id", "username", "first_name", "last_name", "chat_id").
		From("tg_user").
		ToSql()
	if err != nil {
		u.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return nil, err
	}

	rows, err := u.db.Query(sql)
	if err != nil {
		u.logger.Error(fmt.Sprintf("%s : failet exec query %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.Id,
			&user.TgUserID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.ChatId,
		); err != nil {
			u.logger.Error(fmt.Sprintf("%s : failet scan rows %s", op, sql), sl.Err(err))
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
