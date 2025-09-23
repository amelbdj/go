package main

import (
	"fmt"
	"go-projet/function"
)

func main() {
	for {
		function.AfficherMenu()

		var choix int
		fmt.Scanln(&choix)

		if choix == 1 {
			function.AfficherLivre()
		} else if choix == 2 {
			function.AjouterLivre()
		} else if choix == 3 {
			function.SupprimerLivre()
		} else if choix == 4 {
			function.ModifierLivre()
		} else if choix == 5 {
			function.TrierParTitre()
		} else if choix == 6 {
			function.TrierParAnnee()
		} else if choix == 7 {
			fmt.Println("finito")
			return
		} else {
			fmt.Println("Option invalide, réessayez.")
		}

		fmt.Println()
	}
}
