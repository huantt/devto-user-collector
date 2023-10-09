package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
	"user-collector/pkg/forem"
)

type UserRepository struct {
	db        *sql.DB
	tableName string
}

func NewUserRepository(db *sql.DB, tableName string) *UserRepository {
	return &UserRepository{
		db:        db,
		tableName: tableName,
	}
}
func (j *UserRepository) Save(ctx context.Context, user forem.User) error {
	_, err := j.db.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			type_of,
			id,
			username,
			NAME,
			twitter_username,
			github_username,
			summary,
			location,
			website_url,
			joined_at,
			profile_image
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT DO NOTHING
	`, j.tableName),
		user.TypeOf,
		user.Id,
		user.Username,
		user.Name,
		user.TwitterUsername,
		user.GithubUsername,
		user.Summary,
		user.Location,
		user.WebsiteUrl,
		user.JoinedAt,
		user.ProfileImage,
	)
	if err != nil && strings.HasPrefix(err.Error(), "no such table") {
		err := j.createTable()
		if err != nil {
			return err
		}
		return j.Save(ctx, user)
	}
	return err
}

func (j *UserRepository) createTable() error {
	_, err := j.db.Exec(fmt.Sprintf(`
	CREATE TABLE if NOT EXISTS %s (
	    type_of VARCHAR(255),
		id INTEGER,
		username VARCHAR(255),
		NAME VARCHAR(255),
		twitter_username VARCHAR(255),
		github_username VARCHAR(255),
		summary VARCHAR(255),
		location VARCHAR(255),
		website_url VARCHAR(255),
	    joined_at VARCHAR(255),
	    profile_image VARCHAR(255)
	)
`, j.tableName))
	return err
}

func (j *UserRepository) GetLastUserID(ctx context.Context, max int) (int64, error) {
	queryResult, err := j.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT MAX(id)
		FROM %s
		WHERE id <= %d
		ORDER BY id DESC
		LIMIT 1`, j.tableName, max))
	if err != nil {
		if strings.HasPrefix(err.Error(), "no such table") {
			if err := j.createTable(); err != nil {
				return 0, err
			}
			return j.GetLastUserID(ctx, max)
		}
		return 0, err
	}
	defer queryResult.Close()
	var userID *int64
	if !queryResult.Next() {
		return 0, nil
	}
	if err := queryResult.Scan(&userID); err != nil {
		return 0, err
	}
	if userID == nil {
		return 0, nil
	}
	return *userID, nil
}
