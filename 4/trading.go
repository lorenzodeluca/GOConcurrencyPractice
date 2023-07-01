package main

import (
	"math/rand"
	"sync"
	"time"
)

type Quotazione struct {
	prezzo chan float32
	min    float32
	max    float32
	mutex  *sync.Mutex
}

func simulateMarketData(quotazione Quotazione, tradingInCorso *bool) {
	for *tradingInCorso {
		quotazione.mutex.Lock()
		quotRandInRange := (quotazione.min + rand.Float32()*(quotazione.max-quotazione.min))
		quotazione.prezzo <- quotRandInRange
		quotazione.mutex.Unlock()
		time.Sleep(time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) //per avere ad ogni esecuzione una diversa sequenza random

	var tradingInCorso bool = true
	var mutexEU, mutexGU, mutexJU sync.Mutex
	eurusdQuotazione := Quotazione{make(chan float32), 1.0, 1.5, &mutexEU}
	gbpusdQuotazione := Quotazione{make(chan float32), 1.0, 1.5, &mutexGU}
	jpyusdQuotazione := Quotazione{make(chan float32), 1.0, 1.5, &mutexJU}

	go simulateMarketData(eurusdQuotazione, &tradingInCorso)
	go simulateMarketData(gbpusdQuotazione, &tradingInCorso)
	go simulateMarketData(jpyusdQuotazione, &tradingInCorso)
}
