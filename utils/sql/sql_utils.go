package sql

import (
	"database/sql"
	"encoding/json"
	"strings"

	"arion_shot_api/datasources/my_sql/arion_shots_db"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// CloseStmt closes an statement. We are using this function to fix the lint issue
// caused because this function is returning an error object.
func CloseStmt(stmt *sql.Stmt) {
	_ = stmt.Close()
}

// CloseRows closes the Rows object. We are using this function to fix the lint issue
// caused because this function is returning an error object.
func CloseRows(rows *sql.Rows) {
	_ = rows.Close()
}

// CreateStmt creates an statement. We are using this function to reduce the boilerplate
// code.
func CreateStmt(query string) (*sql.Stmt, error) {
	stmt, err := arion_shots_db.Client.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "database error")
	}

	return stmt, nil
}

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) error {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		switch true {
		case strings.Contains(err.Error(), errorNoRows):
			return errors.New("no record matching given information")
		default:
			return errors.Wrap(err, "error parsing database response")
		}
	}

	switch sqlErr.Number {
	case 1062:
		return errors.New("invalid data")
	}

	return errors.Wrap(err, "database error")
}

//NullString is a wrapper around sql.NullString
type NullString sql.NullString

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}
