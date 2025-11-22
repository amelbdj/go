package bdd

import (
	"fmt"
	"game/models"
)

func AddGameInLibrarie(idGame int, idLibrarie int) error {

	// 1️⃣ Vérifier que la librairie existe et récupérer premium
	var isPremium bool
	err := Conn.QueryRow(
		"SELECT is_premium FROM partiel.librarie WHERE id = ?",
		idLibrarie,
	).Scan(&isPremium)

	if err != nil {
		return fmt.Errorf("library not found: %v", err)
	}

	// 2️⃣ Vérifier que le jeu existe
	var gameID int
	err = Conn.QueryRow(
		"SELECT id FROM partiel.game WHERE id = ?",
		idGame,
	).Scan(&gameID)

	if err != nil {
		return fmt.Errorf("game not found: %v", err)
	}

	// 3️⃣ Vérifier que le jeu n'est pas déjà dans la librairie
	var count int
	err = Conn.QueryRow(
		"SELECT COUNT(*) FROM partiel.library_games WHERE library_id = ? AND game_id = ?",
		idLibrarie, idGame,
	).Scan(&count)

	if err != nil {
		return fmt.Errorf("error checking existing association: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("game already in this library")
	}

	// 4️⃣ Vérifier la limite si la librairie n'est pas premium
	err = Conn.QueryRow(
		"SELECT COUNT(*) FROM partiel.library_games WHERE library_id = ?",
		idLibrarie,
	).Scan(&count)

	if err != nil {
		return fmt.Errorf("error checking library game count: %v", err)
	}

	if !isPremium && count >= 5 {
		return fmt.Errorf("non-premium library limit reached (5 games max)")
	}

	// 5️⃣ Ajouter le jeu à la librairie
	_, err = Conn.Exec(
		"INSERT INTO partiel.library_games (library_id, game_id) VALUES (?, ?)",
		idLibrarie, idGame,
	)

	if err != nil {
		return fmt.Errorf("error inserting game into library: %v", err)
	}

	// Tout s'est bien passé
	return nil
}

func DeleteGameInLib(idGame int, idLibrarie int) error {

	// Vérifier si la lib existe

	_, err := Conn.Query("SELECT id FROM partiel.librarie WHERE id = ?", idLibrarie)
	if err != nil {

		return fmt.Errorf("la librarie n'existe pas : %d", idLibrarie)

	}

	//  Vérifier que le jeu existe
	var gameID int
	err = Conn.QueryRow(
		"SELECT id FROM partiel.game WHERE id = ?",
		idGame,
	).Scan(&gameID)

	if err != nil {
		fmt.Println(" Game not found:", err)

	}

	// Suppr le jeu dans la librairie
	_, err = Conn.Exec(
		"DELETE FROM partiel.library_games WHERE library_id = ? AND game_id = ?",
		idLibrarie, idGame,
	)

	if err != nil {
		fmt.Println(" DELETE error:", err)

	}

	fmt.Println(" Game successfully DELETED FROM library")

	return nil

}

func GetGamesInLibrary(idLib int) ([]models.Game, error) {

	// 1️⃣ Vérifier que la librairie existe
	var exists int
	err := Conn.QueryRow(
		"SELECT COUNT(*) FROM partiel.librarie WHERE id = ?",
		idLib,
	).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("error checking library: %v", err)
	}

	if exists == 0 {
		return nil, fmt.Errorf("library not found")
	}

	// 2️⃣ Récupérer les jeux de la librairie
	rows, err := Conn.Query(`
        SELECT partiel.game.id,
               partiel.game.name,
               partiel.game.price
        FROM partiel.game
        JOIN partiel.library_games
            ON partiel.game.id = partiel.library_games.game_id
        WHERE partiel.library_games.library_id = ?`,
		idLib,
	)

	if err != nil {
		return nil, fmt.Errorf("error fetching games: %v", err)
	}
	defer rows.Close()

	var games []models.Game

	// 3️⃣ Parcourir les résultats
	for rows.Next() {
		var game models.Game
		err := rows.Scan(&game.Id, &game.Name, &game.Price)
		if err != nil {
			return nil, fmt.Errorf("error scanning game: %v", err)
		}
		games = append(games, game)
	}

	return games, nil
}
