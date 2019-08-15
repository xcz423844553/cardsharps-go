package main

import (
	"errors"
	"math"
)

type Sharper struct {
}

// func (sharper *Sharper) CalcMACDStat(data []float64, emaRangeShort int, emaRangeLong int, macdRange int) (float64, error) {
// 	//Minimum length of data is PARAM_MACD_RANGE * 2+PARAM_EMA_LONG_RANGE * 2
// 	PARAM_EMA_SHORT_RANGE := emaRangeShort
// 	PARAM_EMA_LONG_RANGE := emaRangeLong
// 	PARAM_MACD_RANGE := macdRange

// 	if len(data) < PARAM_MACD_RANGE*2+PARAM_EMA_LONG_RANGE*2 {
// 		return 0, errors.New("Not enough data to calculate MACD")
// 	}

// 	var emaShort float64 //Short Exponential Moving Average
// 	var emaLong float64  //Long Exponential Moving Average
// 	var diff []float64
// 	var dea float64
// 	var macd float64

// 	var emaShortErr error
// 	var emaLongErr error
// 	var deaErr error

// 	for i := 0; i < PARAM_MACD_RANGE*2; i++ {

// 		// fmt.Printf("MACD #: %d\r\n", i)
// 		emaShort, emaShortErr = sharper.CalcEMAStat(data[i:], PARAM_EMA_SHORT_RANGE)
// 		if emaShortErr != nil {
// 			return macd, emaShortErr
// 		}
// 		// fmt.Printf("SHORT: %.2f\r\n", emaShort)
// 		emaLong, emaLongErr = sharper.CalcEMAStat(data[i:], PARAM_EMA_LONG_RANGE)
// 		if emaLongErr != nil {
// 			return macd, emaLongErr
// 		}
// 		// fmt.Printf("LONG: %.2f\r\n", emaLong)
// 		diff = append(diff, emaShort-emaLong)
// 	}
// 	dea, deaErr = sharper.CalcEMAStat(diff, PARAM_MACD_RANGE)
// 	if deaErr != nil {
// 		return macd, deaErr
// 	}
// 	macd = diff[0] - dea
// 	return macd, nil
// }

// func (sharper *Sharper) CalcEMAStat(data []float64, emaRange int) (float64, error) {
// 	PARAM_EMA_RANGE := emaRange
// 	var ema float64
// 	if len(data) < PARAM_EMA_RANGE*2 || PARAM_EMA_RANGE <= 0 {
// 		return ema, errors.New("Not enough data to calculate EMA")
// 	}
// 	EMA_MULTIPLIER := 2 / float64(PARAM_EMA_RANGE+1)

// 	var ma float64 //Moving Average

// 	//Calculate MA(EMA_RANGE)
// 	for i := PARAM_EMA_RANGE; i < PARAM_EMA_RANGE*2; i++ {
// 		ma += float64(data[i])
// 	}
// 	ma /= float64(PARAM_EMA_RANGE)

// 	//Calculate EMA(EMA_RANGE)
// 	var emaPrev float64
// 	emaPrev = ma
// 	for i := PARAM_EMA_RANGE - 1; i >= 0; i-- {
// 		ema = float64(data[i])*EMA_MULTIPLIER + emaPrev*(1-EMA_MULTIPLIER)
// 		emaPrev = ema
// 	}
// 	return ema, nil
// }

func (sharper *Sharper) CalcBollKcStat(histList []RowStockHist, targetIndex int, crossRange int, lookbackRange int) (float32, float32, float32, float32, float32, float32) {
	PARAM_MA_RANGE := lookbackRange
	PARAM_EMA_PREV_RANGE := lookbackRange
	PARAM_ATR_RANGE := lookbackRange / 2
	PARAM_NUM_SIGMA := 2
	EMA_MULTIPLIER := 2 / float64(PARAM_MA_RANGE+1)

	var ma float64        //Moving Average
	var sd float64        //Standard Deviation
	var ema float64       //Exponential Moving Average
	var atr float64       //Average True Range
	var bollMid float64   //Mid Line of Bollinger Band
	var bollUpper float64 //Upper Bound of Bollinger Band
	var bollLower float64 //Lower Bound of Bollinger Band
	var kcMid float64     //Mid Line of Keltner Channel
	var kcUpper float64   //Upper Bound of Keltner CHannel
	var kcLower float64   //Lower Bound of Keltner Channel

	//Calculate MA(MA_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
		ma += float64((histList[i].MarketHigh + histList[i].MarketLow + 2*histList[i].MarketClose) / 4)
	}
	ma /= float64(PARAM_MA_RANGE)

	//Calculate SIGMA(MA_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
		sdElm := float64((histList[i].MarketHigh+histList[i].MarketLow+2*histList[i].MarketClose)/4) - ma
		sd += math.Pow(sdElm, 2)
	}
	sd = math.Sqrt(sd / float64(PARAM_MA_RANGE))

	//Calculate EMA(MA_RANGE, EMA_PREV_RANGE), including the real time stock price
	var emaPrev float64
	for i := targetIndex + PARAM_MA_RANGE; i < targetIndex+PARAM_MA_RANGE+PARAM_EMA_PREV_RANGE; i++ {
		emaPrev += float64((histList[i].MarketHigh + histList[i].MarketLow + 2*histList[i].MarketClose) / 4)
	}
	emaPrev /= float64(PARAM_EMA_PREV_RANGE)
	for i := targetIndex + PARAM_MA_RANGE - 1; i >= targetIndex; i-- {
		ema = float64((histList[i].MarketHigh+histList[i].MarketLow+2*histList[i].MarketClose)/4)*EMA_MULTIPLIER + emaPrev*(1-EMA_MULTIPLIER)
		emaPrev = ema
	}

	//Calculate ATR(ATR_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_ATR_RANGE; i++ {
		trElm1 := float64(histList[i].MarketHigh - histList[i].MarketLow)
		trElm2 := math.Abs(float64(histList[i].MarketHigh - histList[i+1].MarketClose))
		trElm3 := math.Abs(float64(histList[i].MarketLow - histList[i+1].MarketClose))
		tr := math.Max(math.Max(trElm1, trElm2), trElm3)
		atr += tr
	}
	atr = atr / float64(PARAM_ATR_RANGE)

	//Calculate Parameters of Bollinger Band
	bollMid = ma
	bollUpper = bollMid + float64(PARAM_NUM_SIGMA)*sd
	bollLower = bollMid - float64(PARAM_NUM_SIGMA)*sd
	//Calculate Parameters of Keltner Channel
	kcMid = ema
	kcUpper = kcMid + float64(PARAM_NUM_SIGMA)*atr
	kcLower = kcMid - float64(PARAM_NUM_SIGMA)*atr
	// fmt.Printf("Index:%v Symbol:%s CLOSE: %.2f MA:%.2f SD:%.2f EMA:%.2f ATR:%.2f BMID:%.2f BUP:%.2f BLO:%.2f KMID:%.2f KUP:%.2f KLO:%.2f\r\n", targetIndex, histList[targetIndex].Symbol, histList[targetIndex].MarketClose, ma, sd, ema, atr, bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower)
	return float32(bollMid), float32(bollUpper), float32(bollLower), float32(kcMid), float32(kcUpper), float32(kcLower)
}

func (sharper *Sharper) CalcEMAStat(data []float64, emaRange int) ([]float64, error) {
	if len(data) < emaRange+1 || emaRange <= 0 {
		return nil, errors.New("Not enough data to calculate EMA")
	}
	ema := make([]float64, len(data)-emaRange)
	multiplier := 2 / float64(emaRange+1)

	var emaPrev float64
	for i := len(ema) - 1; i >= 0; i-- {
		//Earliest EMA uses SMA as EMAPrev
		if i == len(ema)-1 {
			//Calculate SMA(EMA_RANGE)
			var sum float64
			for j := len(ema); j < len(data); j++ {
				sum += float64(data[j])
			}
			emaPrev = sum / float64(emaRange)
		}
		//Calculate EMA using EMAPrev
		ema[i] = float64(data[i])*multiplier + emaPrev*(1-multiplier)
		emaPrev = ema[i]
	}
	return ema, nil
}

func (sharper *Sharper) CalcMACDStat(data []float64, emaRangeShort int, emaRangeLong int, macdRange int) ([]float64, []float64, []float64, error) {
	//Minimum length of data is macdRange * 2+emaRangeLong * 2 (if emaRangeLong is larger than emaRangeShort)
	if len(data) < MinInt(macdRange*2+emaRangeLong*2, macdRange*2+emaRangeShort*2) {
		return nil, nil, nil, errors.New("Not enough data to calculate MACD")
	}

	var emaShort []float64 //Short Exponential Moving Average
	var emaLong []float64  //Long Exponential Moving Average
	var diff []float64
	var dea []float64
	var macd []float64

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
	macd = make([]float64, MinInt(len(diff), len(dea)))
	for i := 0; i < MinInt(len(diff), len(dea)); i++ {
		macd[i] = (diff[i] - dea[i]) * 2
	}
	return macd, diff, dea, nil
}
