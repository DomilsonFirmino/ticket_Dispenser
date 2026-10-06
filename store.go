package main

import (
	"sync"
)

var mu sync.RWMutex

// bloqueio para impedir escrita multipla e leitura multipla
var pessoas = map[int]person{
	0: {
		Id:        "0",
		Firstname: "domilson",
		Lastname:  "firmino",
		Age:       10,
	},
}

var pessoaCounter = 1
