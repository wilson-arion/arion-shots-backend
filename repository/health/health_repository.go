package health

import (
	sqlutils "arion_shot_api/internal/utils/sql"
)

const (
	sql = `SELECT true`
)

func Check() (bool, error) {
	stmt, err := sqlutils.CreateStmt(sql) //nolint:sqlclosecheck

	if err != nil {
		return false, err
	}
	defer sqlutils.CloseStmt(stmt)

	var tmp bool

	result := stmt.QueryRow()
	if err := result.Scan(&tmp); err != nil {
		return false, sqlutils.ParseError(err)
	}

	return true, nil
}
