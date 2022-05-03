package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"api.stanible.com/wallet/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := ":" + os.Getenv("SERVER_PORT")
	r := api.Router()
	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
