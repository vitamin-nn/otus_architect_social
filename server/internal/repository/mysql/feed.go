package mysql

import (
	"context"
	"fmt"

	"github.com/vitamin-nn/otus_architect_social/server/internal/db/replication"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

var _ repository.FeedRepo = (*FeedRepo)(nil)

type FeedRepo struct {
	dbPool *replication.DBReplPool
}

func NewFeedRepo(dbPool *replication.DBReplPool) *FeedRepo {
	return &FeedRepo{
		dbPool: dbPool,
	}
}

func (m *FeedRepo) Create(ctx context.Context, feed *repository.Feed) (*repository.Feed, error) {
	db := m.dbPool.GetMaster()
	stmt, err := db.PrepareContext(
		ctx,
		`INSERT INTO user_feed(
			user_id,
			create_at,
			body
		)
		VALUES(?, ?, ?)`,
	)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		feed.UserID,
		feed.CreateAt,
		feed.Body,
	)
	if err != nil {
		specErr := fmt.Errorf("insert error: %v", err)

		return nil, specErr
	}

	feedID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	feed.ID = int(feedID)

	return feed, nil
}

func (m *FeedRepo) GetFeedByUserID(ctx context.Context, userID, limit, offset int) ([]*repository.Feed, error) {
	var result []*repository.Feed

	q := `SELECT
			id,
			user_id,
			create_at,
			body
		FROM user_feed
		WHERE user_id = ?
		ORDER BY id DESC
		LIMIT ?
		OFFSET ?`

	db := m.dbPool.GetSlave()
	rows, err := db.QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var f *repository.Feed
	for rows.Next() {
		f = new(repository.Feed)
		err = rows.Scan(&f.ID,
			&f.UserID,
			&f.CreateAt,
			&f.Body,
		)

		if err != nil {
			return result, err
		}

		result = append(result, f)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
