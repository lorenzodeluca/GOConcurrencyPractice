package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Quotazione struct {
	nome        string
	prezzo      chan float32
	minPrice    float32
	maxPrice    float32
	vendiSopra  float32
	compraSotto float32
}

func simulateMarketData(quotazione Quotazione, tradingInCorso *bool) {
	for *tradingInCorso {
		quotRandInRange := (quotazione.minPrice + rand.Float32()*(quotazione.maxPrice-quotazione.minPrice))
		select {
		case _, ok := <-quotazione.prezzo: // entro in questo caso se nessuno ha preso la quotazione che era attualmente nel canale
			if ok {
				quotazione.prezzo <- quotRandInRange //la riga prima svuota il canale quindi l'inserimento nel canale verrà eseguito sicuramente subito
			} //  ok == false significherebbe canale chiuso
		default: //nessuna valore attualmente nel canale quindi mi limito a caricarci la quotazione aggiornata
			quotazione.prezzo <- quotRandInRange
		}
		time.Sleep(time.Second)
	}

}

func selectPair(eurusd Quotazione, gbpusd Quotazione, jpyusd Quotazione, tradingInCorso *bool, operazioniInCorso *sync.WaitGroup) {
	for *tradingInCorso {
		select {
		case eurusdPrezzo := <-eurusd.prezzo:
			if eurusdPrezzo > eurusd.vendiSopra {
				fmt.Print("Vendita di EUR/USD in corso... ")
				time.Sleep(4 * time.Second)
				fmt.Println(" VENDITA CONFERMATA(Prezzo:", eurusdPrezzo, ")")
			}

		case gbpusdPrezzo := <-gbpusd.prezzo:
			if gbpusdPrezzo < gbpusd.compraSotto {
				fmt.Print("Acquisto di GBP/USD in corso... ")
				time.Sleep(3 * time.Second)
				fmt.Println(" ACQUISTO CONFERMATO(Prezzo:", gbpusdPrezzo, ")")
			}

		case jpyusdPrezzo := <-jpyusd.prezzo:
			if jpyusdPrezzo < jpyusd.compraSotto {
				fmt.Print("Acquisto di JPY/USD in corso... ")
				time.Sleep(3 * time.Second)
				fmt.Println(" ACQUISTO CONFERMATO(Prezzo:", jpyusdPrezzo, ")")
			}
		}
	}
	defer operazioniInCorso.Done()
}

func main() {
	rand.Seed(time.Now().UnixNano()) //per avere ad ogni esecuzione una diversa sequenza random

	var tradingInCorso bool = true
	var operazioniInCorso sync.WaitGroup
	operazioniInCorso.Add(1)

	//999 e -999 sono valori flag(valore che non può sicuramente essere raggiunto per i limiti massimi e minimi delle varie valute)
	eurusdQuotazione := Quotazione{"EUR/USD", make(chan float32), 1.0, 1.5, 1.20, -999}
	gbpusdQuotazione := Quotazione{"GBP/USD", make(chan float32), 1.0, 1.5, 999, 1.35}
	jpyusdQuotazione := Quotazione{"JPY/USD", make(chan float32), 0.006, 0.009, 999, 0.0085}

	go simulateMarketData(eurusdQuotazione, &tradingInCorso)
	go simulateMarketData(gbpusdQuotazione, &tradingInCorso)
	go simulateMarketData(jpyusdQuotazione, &tradingInCorso)

	go selectPair(eurusdQuotazione, gbpusdQuotazione, jpyusdQuotazione, &tradingInCorso, &operazioniInCorso)

	go func() { //routine che gestisce la durata del ciclo di trading
		time.Sleep(time.Minute)
		tradingInCorso = false
	}()

	operazioniInCorso.Wait() //aspetto che l'operazione che è attualmente in esecuzione terminini(se presente)
}
