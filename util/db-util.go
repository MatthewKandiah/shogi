package util

import (
	"os"

	"github.com/MatthewKandiah/shogi/dao"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func InitialiseDb(daos []dao.Dao) error {
	for _, dao := range daos {
		err := dao.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
