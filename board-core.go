package main

import (
	"sync"
)

type Board struct {
	waitGroup sync.WaitGroup
	mutLock   sync.Mutex
}

//SymbolProducer creates a producer of symbols from db_symbol
//Return: outChan - channel to send symbols
func (bd *Board) SymbolProducer(tag string) <-chan string {
	outChan := make(chan string, CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY)

	var symbols []string
	if tag == SYMBOLTAG_ALL {
		rows, selectErr := new(TblSymbol).SelectAllSymbol()
		if selectErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_BOARD, selectErr.Error())
		}
		symbols = append(symbols, rows...)
	} else {
		tagRow := TblSymbolTagRow{
			Tag: tag,
		}
		rows, selectErr := new(TblSymbolTag).SelectSymbolRowByTag(tagRow)
		if selectErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_BOARD, selectErr.Error())
		}
		for _, row := range rows {
			symbols = append(symbols, row.Symbol)
		}
	}
	// symbols = []string{"CANN", "CTST", "SMG", "TRTC", "TSM"}
	// symbols = []string{"AAPL"}
	go func() {
		for _, symbol := range symbols {
			outChan <- symbol
		}
		defer close(outChan)
	}()
	return outChan
}

//SymbolConsumer creates a consumer of symbols
//Param: inChan - channel to receive the symbols
//Param: processFunc - function which will be processed in the consumer
//Return: void
func (bd *Board) SymbolConsumer(inChan <-chan string, execution func(string)) {
	for symbol := range inChan {
		execution(symbol)
	}
	bd.waitGroup.Done()
}

//SymbolGame activates one producer and one consumer, calls a function on all the symbols
//Param: exchange - where the symbols are listed, enum of TRADER_ALL, TRADER_DOW, TRADER_SP500, TRADER_NASDAQ, TRADER_RUSSELL
//Param: processFunc - function which will be processed in the consumer
//Param: callback - callback function executed after the producer and consumer return
//Return: void
func (bd *Board) SymbolGame(tag string, execution func(string), callback func()) {
	symbolChan := bd.SymbolProducer(tag)
	for i := 0; i < CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY; i++ {
		bd.waitGroup.Add(1)
		go bd.SymbolConsumer(symbolChan, execution)
	}
	bd.waitGroup.Wait()
	callback()
	return
}
