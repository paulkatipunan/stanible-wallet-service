package main

import (
	"fmt"
	"log"
	"net/http"

	"api.stanible.com/wallet/router"
)

func main() {
	// connStr := `postgres://admin:Gd0+p2\#Me@@.>iH?9=}Z+M[q9k_<D{@34.142.140.9/dev-wallet?sslmode=disable`
	// db, err := sql.Open("postgres", connStr)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	port := ":8080"
	r := router.Router()
	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
