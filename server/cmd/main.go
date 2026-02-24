package main

import (
	"fmt"
	"linkShorter/internal/router"
	"linkShorter/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env load err")
	}

	st, err := storage.NewStorage(os.Getenv("StorageType"))
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(st)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("ServerPort")), r))

}
