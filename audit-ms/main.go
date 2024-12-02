package main

import (
	"audit-ms/router"
	"log"
	"net/http"
)

func main() {
	r := router.SetupRouter()
	log.Println("Audit microservice is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
