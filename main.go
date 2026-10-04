package main

import (
	"log"
	"net/http"
	"strconv"
)

type JsonResponse struct {
	Sucess  bool   `json:"sucess"`
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /person", AddPerson)
	mux.HandleFunc("GET /person/{id}", GetPerson)
	mux.HandleFunc("GET /persons", GetAllPersons)
	mux.HandleFunc("PATCH /person/{id}", UpdatePerson)
	mux.HandleFunc("DELETE /person/{id}", DeletePerson)

	portNumber := 8081
	log.Printf("Going to listen on port %d\n", portNumber)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(8081), mux))
}
