package bdd

import (
	"fmt"
	"game/models"
)

func GetGames() ([]models.Game, error) {

	var games []models.Game

	rows, err := Conn.Query("SELECT id, name, price FROM partiel.game")

	if err != nil {
		return nil, fmt.Errorf("get games : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {

		var game models.Game

		err := rows.Scan(&game.Id, &game.Name, &game.Price)

		if err != nil {
			return nil, fmt.Errorf("get games : %v", err.Error())
		}
		games = append(games, game)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get games : %v", err.Error())
	}

	return games, nil
}

func CreateGame(game models.Game) error {

	_, err := Conn.Exec("INSERT INTO partiel.game (name, price)VALUES (?,  ?)", game.Name, game.Price)

	if err != nil {
		return fmt.Errorf("CreatGame : %s", err.Error)
	}
	return nil
}

func GetGameByName(name string) ([]models.Game, error) {

	var games []models.Game

	rows, err := Conn.Query("SELECT id, name, price FROM partiel.game WHERE name = ?", name)

	if err != nil {
		return nil, fmt.Errorf("get game by name : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {

		var game models.Game

		err := rows.Scan(&game.Id, &game.Name, &game.Price)

		if err != nil {
			return nil, fmt.Errorf("get game by name : %v", err.Error())
		}
		games = append(games, game)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("get game by name : %v", err.Error())
	}
	return games, nil
}

func GetGameById(id int) ([]models.Game, error) {
	var games []models.Game

	rows, err := Conn.Query("SELECT id, name, price FROM partiel.game WHERE id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("get game by id : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var game models.Game

		err := rows.Scan(&game.Id, &game.Name, &game.Price)

		if err != nil {
			return nil, fmt.Errorf("get game by name : %v", err.Error())
		}

		games = append(games, game)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("get game by id : %v", err.Error())
	}
	return games, nil
}

func ModifyGameById(game models.Game) error {

	// Vérifie si le jeu existe

	_, err := Conn.Query("SELECT id FROM partiel.game WHERE id = ?", game.Id)
	if err != nil {

		return fmt.Errorf("le jeu n'existe pas : %d", game.Id)

	}

	// Met à jour le jeu
	_, err = Conn.Exec(
		"UPDATE partiel.game SET name = ?, price = ? WHERE id = ?",
		game.Name,
		game.Price,
		game.Id,
	)
	if err != nil {
		return fmt.Errorf("mise à jour échouée : %v", err)
	}

	return nil
}
