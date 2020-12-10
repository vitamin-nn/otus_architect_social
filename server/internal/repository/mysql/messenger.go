package mysql

import (
	"context"
	"hash/crc32"
	"strconv"

	"github.com/vitamin-nn/otus_architect_social/server/internal/db/sharding"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

var _ repository.MessengerRepo = (*MessengerRepo)(nil)

type MessengerRepo struct {
	dbPool *sharding.DBShardPool
}

func NewMessengerRepo(dbPool *sharding.DBShardPool) *MessengerRepo {
	return &MessengerRepo{
		dbPool: dbPool,
	}
}

func (m *MessengerRepo) Create(ctx context.Context, message *repository.Message) (*repository.Message, error) {
	dialogID := calculateDialogID(message.SenderID, message.ReceiverID)
	db := m.dbPool.GetShard(dialogID)
	stmt, err := db.PrepareContext(
		ctx,
		`INSERT INTO user_message(
			dialog_id,
			sender_user_id,
			receiver_user_id,
			sent_at,
			body,
			is_read
		)
		VALUES(?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		dialogID,
		message.SenderID,
		message.ReceiverID,
		message.SentAt,
		message.Text,
		message.IsRead,
	)
	if err != nil {
		return nil, err
	}

	messageID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	message.ID = int(messageID)
	message.DialogID = dialogID

	return message, nil
}

func (m *MessengerRepo) GetDialogMessageList(ctx context.Context, user1ID, user2ID, limit, offset int) ([]*repository.Message, error) {
	dialogID := calculateDialogID(user1ID, user2ID)
	db := m.dbPool.GetShard(dialogID)

	var result []*repository.Message

	q := `SELECT
			id,
			dialog_id,
			sender_user_id,
			receiver_user_id,
			sent_at,
			body,
			is_read
		FROM user_message
		WHERE dialog_id = ?
		ORDER BY id DESC
		LIMIT ?
		OFFSET ?`

	rows, err := db.QueryContext(
		ctx,
		q,
		dialogID,
		limit,
		offset,
	)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var ms *repository.Message
	for rows.Next() {
		ms = new(repository.Message)
		err = rows.Scan(&ms.ID,
			&ms.DialogID,
			&ms.SenderID,
			&ms.ReceiverID,
			&ms.SentAt,
			&ms.Text,
			&ms.IsRead,
		)

		if err != nil {
			return result, err
		}

		result = append(result, ms)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func calculateDialogID(senderID, receiverID int) int {
	senderStr := strconv.Itoa(senderID)
	receiverStr := strconv.Itoa(receiverID)
	var s string
	if senderID > receiverID {
		s = receiverStr + "+" + senderStr
	} else {
		s = senderStr + "+" + receiverStr
	}

	return int(crc32.ChecksumIEEE([]byte(s)))
}
