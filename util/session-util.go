package util

import (
	"slices"
	"time"

	"github.com/MatthewKandiah/shogi/constant"
	"github.com/MatthewKandiah/shogi/dao"
)

func HasValidSession(userId string, sessionId string, expiryTime string, sessionsDao dao.SessionsDao) (bool, error) {
	dbSessionIds, err := sessionsDao.GetAll(userId)
	if err != nil {
		return false, err
	}
	matchingIndex := slices.IndexFunc(dbSessionIds, func(sr dao.SessionsRow) bool {
		return sr.SessionId == sessionId
	})
	if matchingIndex == -1 {
		return false, nil
	}
	parsedExpiryTime, err := time.Parse(constant.TIME_FORMAT, expiryTime)
	if err != nil {
		return false, nil
	}
	isValid := time.Now().Compare(parsedExpiryTime) == -1
	return isValid, nil
}
