package dao

import "database/sql"

type SessionsDao struct {
	Db *sql.DB
}

type SessionsRow struct {
	UserId string
	SessionId string
	ExpiryTime string
}

func (d SessionsDao) Create() error {
	_, err := d.Db.Exec("CREATE TABLE sessions (userId TEXT, sessionId TEXT, expiryTime TEXT);")
	return err
}

func (d SessionsDao) Insert(row SessionsRow) error {
	_, err := d.Db.Exec("INSERT INTO sessions (userId, sessionId, expiryTime) VALUES (?, ?, ?)", row.UserId, row.SessionId, row.ExpiryTime)
	return err
}

func (d SessionsDao) Delete(userId string, sessionId string) (sql.Result, error) {
	return d.Db.Exec("DELETE FROM sessions WHERE userId = ? AND sessionId = ?", userId, sessionId)
}
