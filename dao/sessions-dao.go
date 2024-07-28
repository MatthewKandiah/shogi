package dao

import "database/sql"

type SessionsDao struct {
	Db *sql.DB
}

type SessionsRow struct {
	UserId     string
	SessionId  string
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

func (d SessionsDao) GetAll(userId string) ([]SessionsRow, error) {
	rows, err := d.Db.Query("SELECT userId, sessionId, expiryTime FROM sessions WHERE userId = ?", userId)
	if err != nil {
		return nil, err
	}
	result := []SessionsRow{}
	for rows.Next() {
		var dbUserId string
		var sessionId string
		var expiryTime string
		err = rows.Scan(&dbUserId, &sessionId, &expiryTime)
		if err != nil {
			return nil, err
		}
		result = append(result, SessionsRow{UserId: userId, SessionId: sessionId, ExpiryTime: expiryTime})
	}
	return result, nil
}

func (d SessionsDao) Get(sessionId string) (*SessionsRow, error) {
	rows := d.Db.QueryRow("SELECT userId, sessionId, expiryTime FROM sessions WHERE sessionId = ?", sessionId)
	var userId string
	var dbSessionId string
	var expiryTime string
	err := rows.Scan(&userId, &dbSessionId, &expiryTime)
	if err != nil {
		return nil, err
	}
	result := SessionsRow{UserId: userId, SessionId: dbSessionId, ExpiryTime: expiryTime}
	return &result, nil
}
