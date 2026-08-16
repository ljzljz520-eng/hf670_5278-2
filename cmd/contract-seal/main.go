package main

import (
	"log"
	"net/http"

	"contractseal/internal/contract"
	"contractseal/web"
)

func main() {
	service := contract.NewService(contract.NewFixtureFactory())
	log.Fatal(http.ListenAndServe(":8080", web.NewHandler(service)))
}
