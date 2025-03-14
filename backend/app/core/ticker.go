package core

import "strings"

var tickerSet map[string]struct{}

func IsTickerAllowed(ticker string) bool {
	_, ok := tickerSet[strings.ToLower(ticker)]
	return ok
}

func LoadTickers(tickers []string) {
	if tickerSet == nil {
		tickerSet = make(map[string]struct{})
	}
	for _, t := range tickers {
		tickerSet[strings.ToLower(t)] = struct{}{}
	}
}

func GetAllTickers() []string {
	var tickerList []string
	for key := range tickerSet {
		tickerList = append(tickerList, key)
	}
	return tickerList
}
