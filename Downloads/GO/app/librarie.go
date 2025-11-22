package app

import (
	"encoding/json"
	"fmt"
	"game/bdd"
	"game/models"
	"net/http"
	"strconv"
	"unicode"
)

func verifyLibrarieDto(dto models.Librarie) []string {

	var errorMsg []string

	for _, r := range dto.Owner_name {
		if !unicode.IsLetter(r) {
			errorMsg = append(errorMsg, "uniquement des caractères alphabétique")
			break
		}
	}

	if len(dto.Owner_name) < 2 || len(dto.Owner_name) > 100 {
		errorMsg = append(errorMsg, "Name must have a length between 2 and 100")
	}
	if len(dto.Owner_password) < 5 || len(dto.Owner_password) > 100 {
		errorMsg = append(errorMsg, "Name must have a length between 5 and 100")
	}

	return errorMsg

}

func GetAllLibrarie(w http.ResponseWriter, r *http.Request) {

	libraries, err := bdd.GetLibrarie()

	if err != nil {
		http.Error(w, "erreur de récupération des lib", http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(libraries)

	if err != nil {
		http.Error(w, "erreur de conversion", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", response)
}

func CreateLibrarie(w http.ResponseWriter, r *http.Request) {
	var librarieDto models.Librarie

	err := json.NewDecoder(r.Body).Decode(&librarieDto)

	if err != nil {
		http.Error(w,
			"Impossible de décoder un modèle game au format json",
			http.StatusBadRequest)
		return

	}

	err = bdd.CreateLibrarie(librarieDto)
	if err != nil {
		http.Error(w,
			"erreur dans la création d'un jeu",
			http.StatusInternalServerError)
		return
	}
	var errorMsgs = verifyLibrarieDto(librarieDto)

	if len(errorMsgs) > 0 {
		var errFormated, _ = json.Marshal(errorMsgs)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errFormated), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetLibrarieById(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}

	// 2️⃣ Requête à la base
	librarie, err := bdd.GetLibrarieById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(librarie)
}

func DeletedLibrarie(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodDelete {
	// 	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	// 	return
	// }

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}

	err = bdd.DeletedLibrarie(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "librarie suppr")

}

func PatchLibrarieById(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"ID invalide"}`, http.StatusBadRequest)
		return
	}

	var librarieDto models.Librarie
	err = json.NewDecoder(r.Body).Decode(&librarieDto)

	if err != nil {
		http.Error(w,
			"Impossible de décoder un modèle user au format json",
			http.StatusBadRequest)
		return

	}
	librarieDto.Id = id
	err = bdd.PatchLibrarieById(librarieDto)
	if err != nil {
		http.Error(w,
			"erreur dans la MODIF d'un jeu",
			http.StatusInternalServerError)
		return
	}

}

func AddGameInLibrarie(w http.ResponseWriter, r *http.Request) {

	idGameStr := r.PathValue("game_id")
	idGame, err := strconv.Atoi(idGameStr)
	if err != nil {
		http.Error(w, "error : ID du jeu invalide", http.StatusBadRequest)
		return
	}

	idLibStr := r.PathValue("library_id")
	idLib, err := strconv.Atoi(idLibStr)
	if err != nil {
		http.Error(w, "ID de la librairie invalide", http.StatusBadRequest)
		return
	}

	err = bdd.AddGameInLibrarie(idGame, idLib)
	if err != nil {
		http.Error(w, "error appel", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Jeu ajouté à la librairie avec succès")
}

func DeleteGameInLib(w http.ResponseWriter, r *http.Request) {

	// Récupérer l'ID du jeu
	idGameStr := r.PathValue("game_id")
	idGame, err := strconv.Atoi(idGameStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"ID du jeu invalide"}`, http.StatusBadRequest)
		return
	}

	// Récupérer l'ID de la librairie
	idLibStr := r.PathValue("library_id")
	idLib, err := strconv.Atoi(idLibStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "ID de la librairie invalide", http.StatusBadRequest)
		return
	}

	// Supprimer le jeu
	err = bdd.DeleteGameInLib(idGame, idLib)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Jeu supprimé de la librairie avec succès")
}

func GetGamesInLibrary(w http.ResponseWriter, r *http.Request) {

	// Récupérer l'ID de la librairie depuis l'URL
	idLibStr := r.PathValue("library_id")
	idLib, err := strconv.Atoi(idLibStr)
	if err != nil {
		http.Error(w, `{"error":"ID de la librairie invalide"}`, http.StatusBadRequest)
		return
	}

	// Appel de la fonction BDD
	games, err := bdd.GetGamesInLibrary(idLib)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	// Réponse JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}
