package dao

import "database/sql"

type GamesDao struct {
	Db *sql.DB
}

type GamesRow struct {
	GameId             string
	PlayerId1          string
	PlayerId2          string
	GameStatus         string
	TimeRemainingSecs1 int
	TimeRemainingSecs2 int
}

const (
	GameStatusPlayer1Win = "Player1Win"
	GameStatusPlayer2Win = "Player2Win"
	GameStatusInProgress = "InProgress"
)

func (d GamesDao) Create() error {
	_, err := d.Db.Exec(
		"CREATE TABLE games (gameId TEXT, playerId1 TEXT, playerId2 TEXT, gameStatus TEXT, timeRemainingSecs1 INTEGER, timeRemainingSecs2 INTEGER)",
	)
	return err
}

func (d GamesDao) Get(gameId string) (*GamesRow, error) {
	row := d.Db.QueryRow("SELECT playerId1, playerId2, gameStatus, timeRemainingSecs1, timeRemainingSecs2 FROM games WHERE gameId = ?", gameId)
	var (
		playerId1 string
		playerId2 string
		gameStatus string
		timeRemainingSecs1 int
		timeRemainingSecs2 int
	)
	err := row.Scan(&playerId1, &playerId2, &gameStatus, &timeRemainingSecs1, &timeRemainingSecs2)
	if err != nil {
		return nil, err
	}
	result := GamesRow{GameId: gameId, PlayerId1: playerId1, PlayerId2: playerId2, GameStatus: gameStatus, TimeRemainingSecs1: timeRemainingSecs1, TimeRemainingSecs2: timeRemainingSecs2}
	return &result, err
}

func (d GamesDao) GetAll(userId string) ([]GamesRow, error) {
	rows, err := d.Db.Query("SELECT gameId, playerId1, playerId2, gameStatus, timeRemainingSecs1, timeRemainingSecs2 WHERE playerId1 = ? OR playerId2 = ?", userId, userId)
	if err != nil {
		return nil, err
	}
	result := []GamesRow{}
	for rows.Next() {
		var (
			gameId string
			playerId1 string
			playerId2 string
			gameStatus string
			timeRemainingSecs1 int
			timeRemainingSecs2 int
		)
		err := rows.Scan(&gameId, &playerId1, &playerId2, &gameStatus, &timeRemainingSecs1, &timeRemainingSecs2)
		if err != nil {
			return nil, err
		}
		result = append(result, GamesRow{
			GameId: gameId,
			PlayerId1: playerId1,
			PlayerId2: playerId2,
			GameStatus: gameStatus,
			TimeRemainingSecs1: timeRemainingSecs1,
			TimeRemainingSecs2: timeRemainingSecs2,
		})
	}
	return result, nil
}
