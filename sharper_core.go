package main

import (
	"errors"
	"math"
)

//Sharper is a struct containing the functions to calculate statistics
type Sharper2 struct {
}

//CalcBollKcStat calculates the bollinger band (Mid, Upper, Lower) and the Keltner Channel (Mid, Upper, Lower)
func (sharper *Sharper2) CalcBollKcStat(histList []IStockHist, targetIndex int64, crossRange int64, lookbackRange int64) (float32, float32, float32, float32, float32, float32) {
	PARAM_MA_RANGE := lookbackRange
	PARAM_EMA_PREV_RANGE := lookbackRange
	PARAM_ATR_RANGE := lookbackRange / 2
	PARAM_NUM_SIGMA := 2
	EMA_MULTIPLIER := 2 / float32(PARAM_MA_RANGE+1)

	var ma float32        //Moving Average
	var sd float32        //Standard Deviation
	var ema float32       //Exponential Moving Average
	var atr float32       //Average True Range
	var bollMid float32   //Mid Line of Bollinger Band
	var bollUpper float32 //Upper Bound of Bollinger Band
	var bollLower float32 //Lower Bound of Bollinger Band
	var kcMid float32     //Mid Line of Keltner Channel
	var kcUpper float32   //Upper Bound of Keltner Channel
	var kcLower float32   //Lower Bound of Keltner Channel

	//Calculate MA(MA_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
		ma += (histList[i].GetMarketHigh() + histList[i].GetMarketLow() + 2*histList[i].GetMarketClose()) / 4
	}
	ma /= float32(PARAM_MA_RANGE)

	//Calculate SIGMA(MA_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
		sdElm := (histList[i].GetMarketHigh()+histList[i].GetMarketLow()+2*histList[i].GetMarketClose())/4 - ma
		sd += float32(math.Pow(float64(sdElm), 2))
	}
	sd = float32(math.Sqrt(float64(sd / float32(PARAM_MA_RANGE))))

	//Calculate EMA(MA_RANGE, EMA_PREV_RANGE), including the real time stock price
	var emaPrev float32
	for i := targetIndex + PARAM_MA_RANGE; i < targetIndex+PARAM_MA_RANGE+PARAM_EMA_PREV_RANGE; i++ {
		emaPrev += (histList[i].GetMarketHigh() + histList[i].GetMarketLow() + 2*histList[i].GetMarketClose()) / 4
	}
	emaPrev /= float32(PARAM_EMA_PREV_RANGE)
	for i := targetIndex + PARAM_MA_RANGE - 1; i >= targetIndex; i-- {
		ema = (histList[i].GetMarketHigh()+histList[i].GetMarketLow()+2*histList[i].GetMarketClose())/4*EMA_MULTIPLIER + emaPrev*(1-EMA_MULTIPLIER)
		emaPrev = ema
	}

	//Calculate ATR(ATR_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_ATR_RANGE; i++ {
		trElm1 := float64(histList[i].GetMarketHigh() - histList[i].GetMarketLow())
		trElm2 := math.Abs(float64(histList[i].GetMarketHigh() - histList[i+1].GetMarketClose()))
		trElm3 := math.Abs(float64(histList[i].GetMarketLow() - histList[i+1].GetMarketClose()))
		tr := float32(math.Max(math.Max(trElm1, trElm2), trElm3))
		atr += tr
	}
	atr = atr / float32(PARAM_ATR_RANGE)

	//Calculate Parameters of Bollinger Band
	bollMid = ma
	bollUpper = bollMid + float32(PARAM_NUM_SIGMA)*sd
	bollLower = bollMid - float32(PARAM_NUM_SIGMA)*sd
	//Calculate Parameters of Keltner Channel
	kcMid = ema
	kcUpper = kcMid + float32(PARAM_NUM_SIGMA)*atr
	kcLower = kcMid - float32(PARAM_NUM_SIGMA)*atr
	// fmt.Printf("Index:%v Symbol:%s CLOSE: %.2f MA:%.2f SD:%.2f EMA:%.2f ATR:%.2f BMID:%.2f BUP:%.2f BLO:%.2f KMID:%.2f KUP:%.2f KLO:%.2f\r\n", targetIndex, histList[targetIndex].Symbol, histList[targetIndex].MarketClose, ma, sd, ema, atr, bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower)
	return bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower
}

//CalcEMAStat calculates EMA
func (sharper *Sharper2) CalcEMAStat(data []float32, emaRange int) ([]float32, error) {
	if len(data) < emaRange+1 || emaRange <= 0 {
		return nil, errors.New("Not enough data to calculate EMA")
	}
	ema := make([]float32, len(data)-emaRange)
	multiplier := 2 / float32(emaRange+1)

	var emaPrev float32
	for i := len(ema) - 1; i >= 0; i-- {
		//Earliest EMA uses SMA as EMAPrev
		if i == len(ema)-1 {
			//Calculate SMA(EMA_RANGE)
			var sum float32
			for j := len(ema); j < len(data); j++ {
				sum += data[j]
			}
			emaPrev = sum / float32(emaRange)
		}
		//Calculate EMA using EMAPrev
		ema[i] = data[i]*multiplier + emaPrev*(1-multiplier)
		emaPrev = ema[i]
	}
	return ema, nil
}

//CalcMACDStat calculates MACD and returns MACD, DIFF, and DEA
func (sharper *Sharper2) CalcMACDStat(data []float32, emaRangeShort int, emaRangeLong int, macdRange int) ([]float32, []float32, []float32, error) {
	//Minimum length of data is macdRange * 2+emaRangeLong * 2 (if emaRangeLong is larger than emaRangeShort)
	if len(data) < MinInt(macdRange*2+emaRangeLong*2, macdRange*2+emaRangeShort*2) {
		return nil, nil, nil, errors.New("Not enough data to calculate MACD")
	}

	var emaShort []float32 //Short Exponential Moving Average
	var emaLong []float32  //Long Exponential Moving Average
	var diff []float32
	var dea []float32
	var macd []float32

	var emaShortErr error
	var emaLongErr error
	var deaErr error

	emaShort, emaShortErr = sharper.CalcEMAStat(data, emaRangeShort)
	if emaShortErr != nil {
		return macd, diff, dea, emaShortErr
	}
	emaLong, emaLongErr = sharper.CalcEMAStat(data, emaRangeLong)
	if emaLongErr != nil {
		return macd, diff, dea, emaLongErr
	}
	for i := 0; i < MinInt(len(emaLong), len(emaShort)); i++ {
		diff = append(diff, emaShort[i]-emaLong[i])
	}
	dea, deaErr = sharper.CalcEMAStat(diff, macdRange)
	if deaErr != nil {
		return macd, diff, dea, deaErr
	}
	macd = make([]float32, MinInt(len(diff), len(dea)))
	for i := 0; i < MinInt(len(diff), len(dea)); i++ {
		macd[i] = (diff[i] - dea[i]) * 2
	}
	return macd, diff, dea, nil
}
