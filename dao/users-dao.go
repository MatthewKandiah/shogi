package dao

import "database/sql"

type UsersDao struct {
	Db *sql.DB
}

type UsersRow struct {
	Id       string
	UserName string
}

func (d UsersDao) Create() error {
	_, err := d.Db.Exec("CREATE TABLE users (id TEXT, userName TEXT);")
	return err
}

func (d UsersDao) Insert(r UsersRow) error {
	_, err := d.Db.Exec("INSERT INTO users (id, userName) VALUES (?, ?)", r.Id, r.UserName)
	return err
}

func (d UsersDao) Get(userId string) (*UsersRow, error) {
	row := d.Db.QueryRow("SELECT id, userName FROM users WHERE id = ?", userId)
	var id string
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		return nil, err
	}
	result := UsersRow{id, name}
	return &result, nil
}

func (d UsersDao) GetByUserName(un string) (*UsersRow, error) {
	row := d.Db.QueryRow("SELECT id, userName FROM users WHERE userName = ?", un)
	var id string
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		return nil, err
	}
	result := UsersRow{id, name}
	return &result, nil
}
