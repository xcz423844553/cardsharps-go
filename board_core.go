package main

import (
	"sync"
)

//Board includes helper functions of producer and consumer
type Board2 struct {
	waitGroup sync.WaitGroup
	mutLock   sync.Mutex
}

//GetSymbolTagAll returns the tags representing all symbols
func (bd *Board2) GetSymbolTagAll() string {
	return "ALL"
}

//ProducerSymbol creates a producer of symbols from db_symbol with tag
//Return: outChan - channel to send symbols
func (bd *Board2) ProducerSymbol(tag string) (<-chan string, error) {
	outChan := make(chan string, CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY)
	var symbols []string
	if TestMode {
		symbols = append(symbols, "AAPL")
		symbols = append(symbols, "BA")
	} else if tag == bd.GetSymbolTagAll() {
		rows, selectErr := new(DaoSymbol).SelectSymbolAll()
		if selectErr != nil {
			return outChan, selectErr
		}
		for _, row := range rows {
			symbols = append(symbols, row.Symbol)
		}
	} else {
		tagRow := RowSymbolTag{
			Tag: tag,
		}
		rows, selectErr := new(DaoSymbolTag).SelectSymbolByTag(tagRow)
		if selectErr != nil {
			return outChan, selectErr
		}
		for _, row := range rows {
			symbols = append(symbols, row.Symbol)
		}
	}
	go func() {
		for _, symbol := range symbols {
			outChan <- symbol
		}
		defer close(outChan)
	}()
	return outChan, nil
}

//ConsumerSymbol creates a consumer of symbols
//Param: inChan - channel to receive the symbols
//Param: processFunc - function which will be processed in the consumer
func (bd *Board2) ConsumerSymbol(inChan <-chan string, executeFunc func(string)) {
	for symbol := range inChan {
		executeFunc(symbol)
	}
	bd.waitGroup.Done()
}

//StartGame activates one producer and one consumer, calls a function on all the symbols
//Param: tag - tag of the symbols to run through
//Param: executeFunc - function which will be executed in the consumer
//Param: callback - callback function executed after the producer and consumer return
//Return: void
func (bd *Board2) StartGame(tag string, executeFunc func(string), callback func()) error {
	symbolChan, producerSymbolErr := bd.ProducerSymbol(tag)
	if producerSymbolErr != nil {
		return producerSymbolErr
	}
	for i := 0; i < CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY; i++ {
		bd.waitGroup.Add(1)
		go bd.ConsumerSymbol(symbolChan, executeFunc)
	}
	bd.waitGroup.Wait()
	callback()
	return nil
}
