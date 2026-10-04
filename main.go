package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type JsonResponse struct {
	Sucess  bool   `json:"sucess"`
	Message string `json:"message"`
}

type PersonJsonResponse struct {
	JsonResponse
	Data []person `json:"data"`
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Atualizar pessoa")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /person", AddPerson)
	mux.HandleFunc("GET /person/{id}", GetPerson)
	mux.HandleFunc("GET /persons", GetAllPersons)
	// mux.HandleFunc("PATCH /person/update", UpdatePerson)
	mux.HandleFunc("DELETE /person/{id}", DeletePerson)

	portNumber := 8081
	log.Printf("Going to listen on port %d\n", portNumber)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(8081), mux))
}
