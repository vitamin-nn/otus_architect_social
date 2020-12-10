package mysql

import (
	"github.com/go-sql-driver/mysql"
)

const (
	ConstraintViolationCode = 1062
)

func getSpecificError(err error, constraintErr error) error {
	if errMy, ok := err.(*mysql.MySQLError); ok {
		if errMy.Number == ConstraintViolationCode {
			return constraintErr
		}
	}

	return nil
}
