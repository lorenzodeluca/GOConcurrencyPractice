package main

import (
	"fmt"
	"sync"
)

func conta(c byte, s byte, ch chan int, wg *sync.WaitGroup) {
	if c == s { //metto un dato(1) nel canale per ogni ripetizione del carattere trovata
		ch <- 1
	}
	defer wg.Done()
}

func main() {
	stringa := "aaaaaaaaaaaabbbbbbbbcccccddddccccccf"
	var carattere byte = 'c'

	ch := make(chan int, len(stringa)) //buffered channel per tenere il conto delle ripetizioni
	var wg sync.WaitGroup              //wait group per sapere quando tutte le goroutines hanno finito di usare il canale

	for i := 0; i < len(stringa); i++ {
		wg.Add(1)
		go conta(stringa[i], carattere, ch, &wg) //avvio delle goroutines
	}

	wg.Wait() //aspetto che tutte le goroutines abbiano finito
	close(ch) //chiudo il canale

	//conto le ripetizioni annotate dalle go routine nel canale
	cont := len(ch)
	fmt.Printf("%d", cont)
}
