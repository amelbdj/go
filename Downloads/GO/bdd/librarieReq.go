package bdd

import (
	"fmt"
	"game/models"
)

func GetLibrarie() ([]models.Librarie, error) {

	var libraries []models.Librarie

	rows, err := Conn.Query("SELECT id, owner_name, is_premium, creation_year FROM partiel.librarie")

	if err != nil {
		return nil, fmt.Errorf("get libraries : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {

		var librarie models.Librarie

		err := rows.Scan(&librarie.Id, &librarie.Owner_name, &librarie.Is_Premium, &librarie.Creation_year)

		if err != nil {
			return nil, fmt.Errorf("get libraries : %v", err.Error())
		}
		libraries = append(libraries, librarie)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get libraries : %v", err.Error())
	}

	return libraries, nil
}

func CreateLibrarie(librarie models.Librarie) error {
	// Exécution de la requête INSERT
	_, err := Conn.Exec(
		"INSERT INTO partiel.librarie (id,owner_name, is_premium, creation_year, owner_password) VALUES (?, ?, ?, ?, ?)",
		librarie.Id,
		librarie.Owner_name,
		librarie.Is_Premium,
		librarie.Creation_year,
		librarie.Owner_password)

	if err != nil {
		// err.Error est une méthode, il faut l'appeler avec ()
		return fmt.Errorf("Create librarie: %s", err.Error())
	}

	return nil
}

func GetLibrarieById(id int) ([]models.Librarie, error) {
	var libraries []models.Librarie

	rows, err := Conn.Query("SELECT id,owner_name, is_premium, creation_year FROM partiel.librarie WHERE id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("get librarie by id : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var librarie models.Librarie

		err := rows.Scan(&librarie.Id, &librarie.Owner_name, &librarie.Is_Premium, &librarie.Creation_year)

		if err != nil {
			return nil, fmt.Errorf("get librarie by name : %v", err.Error())
		}

		libraries = append(libraries, librarie)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("get game by id : %v", err.Error())
	}
	return libraries, nil
}

func DeletedLibrarie(id int) error {

	_, err := Conn.Query("SELECT id FROM partiel.librarie WHERE id = ?", id)
	if err != nil {

		return fmt.Errorf("la librairie n'existe pas : %d", id)

	}

	_, err = Conn.Exec(
		"DELETE FROM partiel.librarie WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("suppression à échoue : %v", err)
	}

	return nil
}

func PatchLibrarieById(librarie models.Librarie) error {

	_, err := Conn.Query("SELECT id FROM partiel.librarie WHERE id = ?", librarie.Id)
	if err != nil {

		return fmt.Errorf("la librarie n'existe pas : %d", librarie.Id)

	}

	_, err = Conn.Exec(
		"UPDATE partiel.librarie SET is_premium = ? WHERE id = ?", librarie.Is_Premium, librarie.Id)
	if err != nil {
		return fmt.Errorf("mise à jour échouée : %v", err)
	}

	return nil
}
