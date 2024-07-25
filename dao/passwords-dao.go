package dao

import "database/sql"

type PasswordsDao struct {
	Db *sql.DB
}

type PasswordsRow struct {
	UserId string
	Password string
}

func (d PasswordsDao) Create() error {
	_, err := d.Db.Exec("CREATE TABLE passwords (userId TEXT, passwords TEXT);")
	return err
}

func (d PasswordsDao) Insert(row PasswordsRow) error {
	_, err := d.Db.Exec("INSERT INTO passwords (userId, password) VALUES (?, ?)",row.UserId, row.Password)
	return err
}

func (d PasswordsDao) Get(userId string) (*PasswordsRow, error) {
	row := d.Db.QueryRow("SELECT (password) FROM passwords WHERE userId = ?", userId)
	var password string
	err := row.Scan(&password)
	if err != nil {
		return nil, err
	}
	result := PasswordsRow{UserId: userId, Password: password}
	return &result, nil
}
