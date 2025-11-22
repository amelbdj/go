package app

import (
	"encoding/json"
	"fmt"
	"game/bdd"
	"game/models"
	"net/http"
	"strconv"
)

func verifyGameDto(dto models.Game) []string {

	var errorMsg []string

	if len(dto.Name) < 4 || len(dto.Name) > 25 {
		errorMsg = append(errorMsg, "Name must have a length between 5 and 50")
	}

	if dto.Price <= 0 {
		errorMsg = append(errorMsg, "Price must not be negative")
	}

	return errorMsg

}

func GetAllGames(w http.ResponseWriter, r *http.Request) {

	games, err := bdd.GetGames()

	if err != nil {
		http.Error(w, "erreur de récupération des jeux", http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(games)

	if err != nil {
		http.Error(w, "erreur de conversion", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", response)
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	var gameDto models.Game

	err := json.NewDecoder(r.Body).Decode(&gameDto)

	if err != nil {
		http.Error(w,
			"Impossible de décoder un modèle game au format json",
			http.StatusBadRequest)
		return

	}

	err = bdd.CreateGame(gameDto)
	if err != nil {
		http.Error(w,
			"erreur dans la création d'un jeu",
			http.StatusInternalServerError)
		return
	}
	var errorMsgs = verifyGameDto(gameDto)

	if len(errorMsgs) > 0 {
		var errFormated, _ = json.Marshal(errorMsgs)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errFormated), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetGameById(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}

	// 2️⃣ Requête à la base
	game, err := bdd.GetGameById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func GetGameByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	game, err := bdd.GetGameByName(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)

}

func ModifyGameById(w http.ResponseWriter, r *http.Request) {

	var gameDto models.Game
	err := json.NewDecoder(r.Body).Decode(&gameDto)

	if err != nil {
		http.Error(w,
			"Impossible de décoder un modèle user au format json",
			http.StatusBadRequest)
		return

	}

	err = bdd.ModifyGameById(gameDto)
	if err != nil {
		http.Error(w,
			"erreur dans la MODIF d'un jeu",
			http.StatusInternalServerError)
		return
	}
	var errorMsgs = verifyGameDto(gameDto)

	if len(errorMsgs) > 0 {
		var errFormated, _ = json.Marshal(errorMsgs)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errFormated), http.StatusBadRequest)
		return

	}

}

func DeletedGame(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}

	err = bdd.DeletedGame(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "jeu suppr")

}
