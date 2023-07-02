package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Torta struct {
	nome     string
	cotta    bool //false=torta non cotta, true=torta cotta
	guarnita bool //false=torta non guarnita
	decorata bool //false=torta non decorata
}

func main() {
	var torteDaFare int = 5
	torteCrudeDaPreparare := make(chan Torta, torteDaFare)
	tortePronte := make(chan Torta, torteDaFare)

	spaziPrimoPasticcere := make(chan Torta, 2)   //spazi per torte cotte
	spaziSecondoPasticcere := make(chan Torta, 2) //spazi per torte guarnite

	var pasticceriAlLavoro sync.WaitGroup //uso il waitGroup per mantenere la main goroutine attiva finchè le altre non finiscono
	pasticceriAlLavoro.Add(3)

	for i := 0; i < torteDaFare; i++ { //creo le torte che dovranno essere preparate
		torteCrudeDaPreparare <- Torta{"Torta#" + strconv.Itoa(i), false, false, false}
	}

	go func() { //goroutine che simula il 1* pasticcere, cucina
		fmt.Println("1 pasticcere: ho iniziato a lavorare")
		torteLavorate := 0
		for torteLavorate < torteDaFare {
			if len(spaziPrimoPasticcere) != cap(spaziPrimoPasticcere) { //controllo se c'è uno spazio libero per poi metterci la torta cotta
				tortaDaCuocere := <-torteCrudeDaPreparare
				fmt.Println("1 pasticcere: ", tortaDaCuocere, " presa in carico")
				time.Sleep(time.Second)
				tortaDaCuocere.cotta = true
				fmt.Println("1 pasticcere:", tortaDaCuocere, " cotta")
				spaziPrimoPasticcere <- tortaDaCuocere
				torteLavorate = torteLavorate + 1
			} else {
				time.Sleep(time.Second) // se non c'è spazio per mettere una torta cotta aspetto un secondo prima di controllare di nuovo
			}
		}
		fmt.Println("1 pasticcere: ho finito di lavorare")
		defer pasticceriAlLavoro.Done()
	}()

	go func() { //goroutine che simula il 2* pasticcere, guarnisce le torte
		fmt.Println("2 pasticcere: ho iniziato a lavorare")
		torteLavorate := 0
		for torteLavorate < torteDaFare {
			if len(spaziSecondoPasticcere) != cap(spaziSecondoPasticcere) { //controllo se c'è uno spazio libero per poi metterci la torta guarnita
				tortaDaGuarnire := <-spaziPrimoPasticcere
				fmt.Println("2 pasticcere: ", tortaDaGuarnire, " presa in carico")
				time.Sleep(4 * time.Second)
				tortaDaGuarnire.guarnita = true
				fmt.Println("2 pasticcere:", tortaDaGuarnire, " guarnita")
				spaziSecondoPasticcere <- tortaDaGuarnire
				torteLavorate = torteLavorate + 1
			} else {
				time.Sleep(time.Second) // se non c'è spazio per mettere una torta guarnita aspetto un secondo prima di controllare di nuovo
			}
		}
		fmt.Println("2 pasticcere: ho finito di lavorare")
		defer pasticceriAlLavoro.Done()
	}()

	go func() { //goroutine che simula il 3* pasticcere, decora le torte guarnite
		fmt.Println("3 pasticcere: ho iniziato a lavorare")
		torteLavorate := 0
		for torteLavorate < torteDaFare {
			tortaDaDecorare := <-spaziSecondoPasticcere
			fmt.Println("3 pasticcere: ", tortaDaDecorare, " presa in carico")
			time.Sleep(8 * time.Second)
			tortaDaDecorare.decorata = true
			fmt.Println("3 pasticcere:", tortaDaDecorare, " decorata")
			tortePronte <- tortaDaDecorare
			torteLavorate = torteLavorate + 1
		}
		fmt.Println("3 pasticcere: ho finito di lavorare")
		defer pasticceriAlLavoro.Done()
	}()

	pasticceriAlLavoro.Wait() //attendo che le goroutine dei tre pasticceri finiscano la loro esecuzione
	fmt.Println("Tutte le torte sono pronte e tutti i pasticceri hanno finito di lavorare!")
}
