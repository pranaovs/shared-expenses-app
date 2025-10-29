package db

import (
	"context"
	"errors"
	"time"

	"shared-expenses-app/models"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateGroup(ctx context.Context, pool *pgxpool.Pool, name, description, userID string) (string, error) {
	var groupID string

	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO groups (group_name, description, created_by, created_at)
			 VALUES ($1, $2, $3, $4)
			 RETURNING group_id`,
		name, description, userID, time.Now(),
	).Scan(&groupID)
	if err != nil {
		return "", err
	}
	// Return the new user's ID
	return groupID, nil
}

func GetGroup(ctx context.Context, pool *pgxpool.Pool, groupID string) (models.Group, error) {
	var group models.Group
	err := pool.QueryRow(
		context.Background(),
		`SELECT name, description, created_by, extract(epoch from created_at)::bigint
		FROM groups
		WHERE group_id = $1`,
	).Scan(&group.Name, &group.Description, &group.CreatedBy, &group.CreatedAt)
	if err == pgx.ErrNoRows {
		return models.Group{}, errors.New("group not found")
	}
	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}
