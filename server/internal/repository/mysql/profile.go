package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vitamin-nn/otus_architect_social/server/internal/db/replication"
	outErr "github.com/vitamin-nn/otus_architect_social/server/internal/error"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

var _ repository.ProfileRepo = (*ProfileRepo)(nil)

type ProfileRepo struct {
	dbPool *replication.DBReplPool
}

func NewProfileRepo(dbPool *replication.DBReplPool) *ProfileRepo {
	return &ProfileRepo{
		dbPool: dbPool,
	}
}

func (m *ProfileRepo) CreateProfile(ctx context.Context, profile *repository.Profile) (*repository.Profile, error) {
	db := m.dbPool.GetMaster()
	stmt, err := db.PrepareContext(
		ctx,
		`INSERT INTO user_profile(
			email,
			password_hash,
			first_name,
			last_name,
			birthdate,
			sex,
			interest_list,
			city
		)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		profile.Email,
		profile.PasswordHash,
		profile.FirstName,
		profile.LastName,
		profile.Birth,
		profile.Sex,
		profile.Interest,
		profile.City,
	)
	if err != nil {
		specErr := getSpecificError(err, outErr.ErrUserAlreadyExists)
		if specErr == nil {
			specErr = fmt.Errorf("insert error: %v", err)
		}

		return nil, specErr
	}

	profileID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	profile.ID = int(profileID)

	return profile, nil
}

func (m *ProfileRepo) UpdateProfile(ctx context.Context, profile *repository.Profile) (*repository.Profile, error) {
	return nil, nil
}

func (m *ProfileRepo) AddFriend(ctx context.Context, profileID1, profileID2 int) error {
	db := m.dbPool.GetMaster()
	stmt, err := db.PrepareContext(
		ctx,
		`INSERT INTO user_friends(
			user_id1,
			user_id2
		)
		VALUES(?, ?), (?, ?)`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		profileID1,
		profileID2,
		profileID2,
		profileID1,
	)
	if err != nil {
		specErr := getSpecificError(err, outErr.ErrFriendAlreadyExists)
		if specErr == nil {
			specErr = fmt.Errorf("insert error: %v", err)
		}

		return specErr
	}

	return nil
}

func (m *ProfileRepo) RemoveFriend(ctx context.Context, profileID1, profileID2 int) error {
	db := m.dbPool.GetMaster()
	stmt, err := db.PrepareContext(
		ctx,
		`DELETE FROM user_friends
		WHERE
		(user_id1 = ? AND user_id2 = ?)
		OR (user_id1 = ? AND user_id2 = ?)`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		profileID1,
		profileID2,
		profileID2,
		profileID1,
	)
	if err != nil {
		specErr := getSpecificError(err, nil)
		if specErr == nil {
			specErr = fmt.Errorf("insert error: %v", err)
		}

		return specErr
	}

	return nil
}

func (m *ProfileRepo) GetProfileByID(ctx context.Context, pID int) (*repository.Profile, error) {
	return m.getRowByCondition(ctx, "id=?", pID)
}

func (m *ProfileRepo) GetProfileByEmail(ctx context.Context, email string) (*repository.Profile, error) {
	return m.getRowByCondition(ctx, "email=?", email)
}

func (m *ProfileRepo) getRowByCondition(ctx context.Context, cond string, val interface{}) (*repository.Profile, error) {
	db := m.dbPool.GetSlave()
	q := fmt.Sprintf(`SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			birthdate,
			sex,
			interest_list,
			city
		FROM user_profile
		WHERE %s LIMIT 1`,
		cond,
	)

	stmt, err := db.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	p := new(repository.Profile)
	err = stmt.QueryRowContext(
		ctx,
		val,
	).Scan(&p.ID,
		&p.Email,
		&p.PasswordHash,
		&p.FirstName,
		&p.LastName,
		&p.Birth,
		&p.Sex,
		&p.Interest,
		&p.City,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		specErr := getSpecificError(err, nil)
		if specErr == nil {
			specErr = fmt.Errorf("sql error: %v", err)
		}

		return nil, specErr
	}

	return p, nil
}

func (m *ProfileRepo) GetProfileList(ctx context.Context, limit, offset int) ([]*repository.Profile, error) {
	var result []*repository.Profile

	q := `SELECT
			id,
			email,
			first_name,
			last_name,
			birthdate,
			sex,
			interest_list,
			city
		FROM user_profile
		ORDER BY id DESC
		LIMIT ?
		OFFSET ?`

	db := m.dbPool.GetSlave()
	rows, err := db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var p *repository.Profile
	for rows.Next() {
		p = new(repository.Profile)
		err = rows.Scan(&p.ID,
			&p.Email,
			&p.FirstName,
			&p.LastName,
			&p.Birth,
			&p.Sex,
			&p.Interest,
			&p.City,
		)

		if err != nil {
			return result, err
		}

		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (m *ProfileRepo) GetFriendsProfileList(ctx context.Context, profileID, limit, offset int) ([]*repository.Profile, error) {
	var result []*repository.Profile

	q := `SELECT
			id,
			email,
			first_name,
			last_name,
			birthdate,
			sex,
			interest_list,
			city
		FROM user_profile
		WHERE id IN (SELECT user_id2 FROM user_friends WHERE user_id1 = ?)
		ORDER BY id DESC
		LIMIT ?
		OFFSET ?`

	db := m.dbPool.GetSlave()
	rows, err := db.QueryContext(ctx, q, profileID, limit, offset)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var p *repository.Profile
	for rows.Next() {
		p = new(repository.Profile)
		err = rows.Scan(&p.ID,
			&p.Email,
			&p.FirstName,
			&p.LastName,
			&p.Birth,
			&p.Sex,
			&p.Interest,
			&p.City,
		)

		if err != nil {
			return result, err
		}

		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (m *ProfileRepo) GetProfileListByNameFilter(ctx context.Context, fName, sName string, limit, offset int) ([]*repository.Profile, error) {
	var result []*repository.Profile

	q := `SELECT
			id,
			email,
			first_name,
			last_name,
			birthdate,
			sex,
			interest_list,
			city
		FROM user_profile FORCE INDEX (f_l_name_idx)
		WHERE first_name LIKE ?
		AND last_name LIKE ?
		ORDER BY id
		LIMIT ?
		OFFSET ?`

	db := m.dbPool.GetSlave()
	rows, err := db.QueryContext(ctx, q, fName, sName, limit, offset)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var p *repository.Profile
	for rows.Next() {
		p = new(repository.Profile)
		err = rows.Scan(&p.ID,
			&p.Email,
			&p.FirstName,
			&p.LastName,
			&p.Birth,
			&p.Sex,
			&p.Interest,
			&p.City,
		)

		if err != nil {
			return result, err
		}

		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
