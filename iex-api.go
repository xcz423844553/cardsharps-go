package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type IexApiManager struct {
}

func (manager *IexApiManager) GetStockChartBySymbolAndRange(symbol string, chartRange string) ([]IexChart, error) {
	var charts []IexChart
	url := manager.GetChartUrlBySymbolAndRange(symbol, chartRange)
	resp, connError := http.Get(url)
	if connError != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER, connError.Error())
		return charts, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER, parseError.Error())
		return charts, parseError
	}
	if jsonError := json.Unmarshal(body, &charts); jsonError != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER, jsonError.Error())
		return charts, parseError
	}
	return charts, nil
}

func (manager *IexApiManager) GetChartUrlBySymbolAndRange(symbol string, chartRange string) string {
	return URL_IEX_CHART_PART1 + symbol + URL_IEX_CHART_PART2 + chartRange
}

func (manager *IexApiManager) WriteIexChartToCsvFile(symbol string, charts []IexChart) error {
	file, fileCreateErr := os.Create("./charts/" + symbol + "_" + strconv.Itoa(GetTimeInYYYYMMDD()) + ".csv")
	if fileCreateErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER, fileCreateErr.Error())
		return fileCreateErr
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"Date", "Minute", "High", "Low",
		"Average", "Volume", "Notional", "NumberOfTrades", "MarketHigh", "MarketLow",
		"MarketAverage", "MarketVolume", "MarketNotional", "MarketNumberOfTrades", "Open",
		"Close", "MarketOpen", "MarketClose", "ChangeOverTime", "MarketChangeOverTime"}
	writeHeaderErr := writer.Write(header)
	if writeHeaderErr != nil {
		return writeHeaderErr
	}
	for _, chart := range charts {
		str := fmt.Sprintf("%v", chart)
		strs := strings.Fields(str[1:(len(str) - 1)])
		writeRowErr := writer.Write(strs)
		if writeRowErr != nil {
			return writeRowErr
		}
	}
	return nil
}
