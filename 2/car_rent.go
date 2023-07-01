package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cliente struct {
	nome string
}
type Veicolo struct {
	tipo string
}

func noleggia(c Cliente, veicoli [3]Veicolo, noleggi map[Cliente]Veicolo, mapMutex *sync.Mutex, wg *sync.WaitGroup) {
	mapMutex.Lock()
	noleggi[c] = veicoli[rand.Intn(3)]
	fmt.Println("Il cliente ", c.nome, "ha noleggiato il veicolo ", noleggi[c])
	mapMutex.Unlock()
	defer wg.Done()
}

func stampa(noleggi map[Cliente]Veicolo, veicoli [3]Veicolo) {
	fmt.Println("Veicoli noleggiati:")
	var noleggiCont = make(map[Veicolo]int) //le variabili intere vengono automaticamente inizializzate a 0
	for _, v := range noleggi {
		noleggiCont[v] = noleggiCont[v] + 1
	}

	for i := range veicoli {
		fmt.Println(veicoli[i].tipo, ": ", noleggiCont[veicoli[i]])
	}
}
func main() {
	rand.Seed(time.Now().UnixNano()) //per avere ad ogni esecuzione una diversa sequenza random
	clienti := [10]Cliente{{"Alessandro Harrington"}, {"Karley Hancock"}, {"Mariam Sawyer"}, {"Jax Kidd"}, {"Jocelyn Parks"}, {"Harry Potter"}, {"Kaylin Franklin"}, {"Aileen Waters"}, {"Jonathan Pollard"}, {"Rylie Evans"}}
	veicoli := [3]Veicolo{{"Berlina"}, {"SUV"}, {"Station Wagon"}}
	var noleggi = make(map[Cliente]Veicolo)
	var wg sync.WaitGroup
	var mapMutex sync.Mutex

	for _, c := range clienti {
		wg.Add(1)
		go noleggia(c, veicoli, noleggi, &mapMutex, &wg)
	}

	wg.Wait()

	stampa(noleggi, veicoli)
}
