package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	outErr "github.com/vitamin-nn/otus_architect_social/server/internal/error"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

const (
	ConstraintViolationCode = 1062
)

var _ repository.ProfileRepo = (*MySQL)(nil)

type MySQL struct {
	db *sql.DB
}

func NewProfileRepo(db *sql.DB) *MySQL {
	return &MySQL{
		db: db,
	}
}

func (m *MySQL) CreateProfile(ctx context.Context, profile *repository.Profile) (*repository.Profile, error) {
	stmt, err := m.db.PrepareContext(
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

func (m *MySQL) UpdateProfile(ctx context.Context, profile *repository.Profile) (*repository.Profile, error) {
	return nil, nil
}

func (m *MySQL) GetProfileByID(ctx context.Context, pID int) (*repository.Profile, error) {
	stmt, err := m.db.PrepareContext(
		ctx,
		`SELECT
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
		WHERE id=? LIMIT 1`,
	)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	p := new(repository.Profile)
	err = stmt.QueryRowContext(
		ctx,
		pID,
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

func (m *MySQL) AddFriend(ctx context.Context, profileID1, profileID2 int) error {
	stmt, err := m.db.PrepareContext(
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

func (m *MySQL) RemoveFriend(ctx context.Context, profileID1, profileID2 int) error {
	stmt, err := m.db.PrepareContext(
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

func (m *MySQL) GetProfileByEmail(ctx context.Context, email string) (*repository.Profile, error) {
	stmt, err := m.db.PrepareContext(
		ctx,
		`SELECT
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
		WHERE email=? LIMIT 1`,
	)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	p := new(repository.Profile)
	err = stmt.QueryRowContext(
		ctx,
		email,
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

func (m *MySQL) GetProfileList(ctx context.Context, limit, offset int) ([]*repository.Profile, error) {
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

	rows, err := m.db.QueryContext(ctx, q, limit, offset)

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

func (m *MySQL) GetFriendsProfileList(ctx context.Context, profileID, limit, offset int) ([]*repository.Profile, error) {
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

	rows, err := m.db.QueryContext(ctx, q, profileID, limit, offset)

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

func (m *MySQL) GetProfileListByNameFilter(ctx context.Context, fName, sName string, limit, offset int) ([]*repository.Profile, error) {
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

	rows, err := m.db.QueryContext(ctx, q, fName, sName, limit, offset)

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

func getSpecificError(err error, constraintErr error) error {
	if errMy, ok := err.(*mysql.MySQLError); ok {
		if errMy.Number == ConstraintViolationCode {
			return constraintErr
			//return outErr.ErrUserAlreadyExists
		}
	}

	return nil
}
