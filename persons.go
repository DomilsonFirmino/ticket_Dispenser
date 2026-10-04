package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type person struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

type PersonJsonResponse struct {
	JsonResponse
	Data []person `json:"data"`
}

type PersonsSingleResponse struct {
	JsonResponse
	Data PersonsResponse `json:"data"`
}

type PersonsMultipleResponse struct {
	JsonResponse
	Data []PersonsResponse `json:"data"`
}

type PersonsResponse struct {
	Type       string
	Id         string
	Attributes person
}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	// bloqueio de acesso, escrita multipla
	mu.Lock()
	defer mu.Unlock()

	contentType := r.Header.Get("Content-type")

	var personS person

	if !strings.HasPrefix(contentType, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Unsupported Media Type: Expected application/json",
		})
		return
	}
	err := json.NewDecoder(r.Body).Decode(&personS)
	// Validacoes mais severas estariao aqui
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: err.Error(),
		})
		return
	}
	switch {
	case personS.Firstname == "":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Firstname is mandatory",
		})
		return
	case personS.Lastname == "":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Lastname is mandatory",
		})
		return
	}

	personS.Id = strconv.Itoa(pessoaCounter)
	// ter de ser requestJson, os valores devem ser a validados,
	pessoas[pessoaCounter] = personS
	pessoaCounter = pessoaCounter + 1

	response := PersonsSingleResponse{
		Sucess:  true,
		Message: "Utilizador criado com sucesso",
		Data: PersonsResponse{
			Type:       "person",
			Id:         strconv.Itoa(pessoaCounter - 1),
			Attributes: personS,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	// bloqueio de acesso, escrita multipla
	mu.Lock()
	defer mu.Unlock()

	//validar o id se existe, se nao existe devolver erro
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "invalid id" + err.Error(),
		})
		return
	}

	PersonF, ok := pessoas[id]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Not found",
		})
		return
	}

	//validar se o content type é json
	contentType := r.Header.Get("Content-type")

	if !strings.HasPrefix(contentType, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Unsupported Media Type: Expected application/json",
		})
		return
	}

	var personS person
	// validar se o body é json, se nao for devolver erro

	decodeErr := json.NewDecoder(r.Body).Decode(&personS)

	fmt.Println("personS dados recebidos: ", personS, "PersonF dados encontrados: ", PersonF)

	if decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: decodeErr.Error(),
		})
		return
	}

	// atualizar o utilizador com os novos dados
	var changeCounter int = 0

	if personS.Firstname != "" && personS.Firstname != PersonF.Firstname {
		PersonF.Firstname = personS.Firstname
	} else if personS.Firstname == PersonF.Firstname {
		changeCounter++
	}

	if personS.Lastname != "" && personS.Lastname != PersonF.Lastname {
		PersonF.Lastname = personS.Lastname
	} else if personS.Lastname == PersonF.Lastname {
		changeCounter++
	}

	if personS.Age != 0 && personS.Age != PersonF.Age {
		PersonF.Age = personS.Age
	} else if personS.Age == PersonF.Age {
		changeCounter++
	}

	if changeCounter == 3 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "No changes made to the user",
		})
		return
	}
	// atualizar no mapa
	pessoas[id] = PersonF

	// devolver o utilizador atualizado
	response := PersonsSingleResponse{
		Sucess:  true,
		Message: "Utilizador Atualizado com sucessso",
		Data: PersonsResponse{
			Type:       "person",
			Id:         strconv.Itoa(id),
			Attributes: PersonF,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	// garantir que ele chegou na rota primeiro

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "invalid id" + err.Error(),
		})
		return
	}

	PersonF, ok := pessoas[id]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Not found",
		})
		return
	}

	response := PersonsSingleResponse{
		Sucess:  true,
		Message: "Utilizador recuperado com sucessso",
		Data: PersonsResponse{
			Type:       "person",
			Id:         strconv.Itoa(id),
			Attributes: PersonF,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAllPersons(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")

	var persons = []PersonsResponse{}

	for index, v := range pessoas {
		currentPerson := PersonsResponse{
			Type:       "person",
			Id:         strconv.Itoa(index),
			Attributes: v,
		}
		persons = append(persons, currentPerson)
	}

	response := PersonsMultipleResponse{
		Sucess:  true,
		Message: "Utilizdores encontrados com sucesso",
		Data:    persons,
	}

	json.NewEncoder(w).Encode(response)

	// Bateram na rota, devolver todos os utilizadores possiveis, sem validações no momento
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// garantir que ele chegou na rota primeiro

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "invalid id" + err.Error(),
		})
		return
	}

	PersonF, ok := pessoas[id]

	if !ok {

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonResponse{
			Sucess:  false,
			Message: "Utilizador não existente",
		})
		return
	}

	var FoundPerson []person
	FoundPerson = append(FoundPerson, PersonF)

	delete(pessoas, id)
	response := PersonJsonResponse{
		Sucess:  true,
		Message: "Utilizadore removido com sucesso",
		Data:    FoundPerson,
	}

	json.NewEncoder(w).Encode(response)
}
