package main

import (
	"fmt"
	"game/app"
	"game/bdd"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	err := bdd.Conn.Ping()

	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "ping à la bdd")
}

func main() {
	bdd.Conn = bdd.NewDB()
	http.HandleFunc("GET /{$}", Health)

	http.HandleFunc("GET /games/", app.GetAllGames)
	http.HandleFunc("POST /games/", app.CreateGame)
	http.HandleFunc("GET /games/{id}/{$}", app.GetGameById)
	http.HandleFunc("DELETE /games/{id}", app.DeletedGame)
	http.HandleFunc("PUT /games/{id}", app.ModifyGameById)
	fmt.Println("test de : http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
