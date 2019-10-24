package main

import (
	"errors"
	"math"
)

//Sharper is a struct containing the functions to calculate statistics
type Sharper2 struct {
}

//CalcBollKcStat calculates the bollinger band (Mid, Upper, Lower) and the Keltner Channel (Mid, Upper, Lower)
// func (sharper *Sharper2) CalcBollKcStat(histList []RowStockHist2, targetIndex int64, crossRange int64, lookbackRange int64) (float32, float32, float32, float32, float32, float32) {
// 	PARAM_MA_RANGE := lookbackRange
// 	PARAM_EMA_PREV_RANGE := lookbackRange
// 	PARAM_ATR_RANGE := lookbackRange / 2
// 	PARAM_NUM_SIGMA := 2
// 	EMA_MULTIPLIER := 2 / float32(PARAM_MA_RANGE+1)

// 	var ma float32        //Moving Average
// 	var sd float32        //Standard Deviation
// 	var ema float32       //Exponential Moving Average
// 	var atr float32       //Average True Range
// 	var bollMid float32   //Mid Line of Bollinger Band
// 	var bollUpper float32 //Upper Bound of Bollinger Band
// 	var bollLower float32 //Lower Bound of Bollinger Band
// 	var kcMid float32     //Mid Line of Keltner Channel
// 	var kcUpper float32   //Upper Bound of Keltner Channel
// 	var kcLower float32   //Lower Bound of Keltner Channel

// 	//Calculate MA(MA_RANGE), including the real time stock price
// 	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
// 		ma += (histList[i].GetMarketHigh() + histList[i].GetMarketLow() + 2*histList[i].GetMarketClose()) / 4
// 	}
// 	ma /= float32(PARAM_MA_RANGE)

// 	//Calculate SIGMA(MA_RANGE), including the real time stock price
// 	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
// 		sdElm := (histList[i].GetMarketHigh()+histList[i].GetMarketLow()+2*histList[i].GetMarketClose())/4 - ma
// 		sd += float32(math.Pow(float64(sdElm), 2))
// 	}
// 	sd = float32(math.Sqrt(float64(sd / float32(PARAM_MA_RANGE))))

// 	//Calculate EMA(MA_RANGE, EMA_PREV_RANGE), including the real time stock price
// 	var emaPrev float32
// 	for i := targetIndex + PARAM_MA_RANGE; i < targetIndex+PARAM_MA_RANGE+PARAM_EMA_PREV_RANGE; i++ {
// 		emaPrev += (histList[i].GetMarketHigh() + histList[i].GetMarketLow() + 2*histList[i].GetMarketClose()) / 4
// 	}
// 	emaPrev /= float32(PARAM_EMA_PREV_RANGE)
// 	for i := targetIndex + PARAM_MA_RANGE - 1; i >= targetIndex; i-- {
// 		ema = (histList[i].GetMarketHigh()+histList[i].GetMarketLow()+2*histList[i].GetMarketClose())/4*EMA_MULTIPLIER + emaPrev*(1-EMA_MULTIPLIER)
// 		emaPrev = ema
// 	}

// 	//Calculate ATR(ATR_RANGE), including the real time stock price
// 	for i := targetIndex; i < targetIndex+PARAM_ATR_RANGE; i++ {
// 		trElm1 := float64(histList[i].GetMarketHigh() - histList[i].GetMarketLow())
// 		trElm2 := math.Abs(float64(histList[i].GetMarketHigh() - histList[i+1].GetMarketClose()))
// 		trElm3 := math.Abs(float64(histList[i].GetMarketLow() - histList[i+1].GetMarketClose()))
// 		tr := float32(math.Max(math.Max(trElm1, trElm2), trElm3))
// 		atr += tr
// 	}
// 	atr = atr / float32(PARAM_ATR_RANGE)

// 	//Calculate Parameters of Bollinger Band
// 	bollMid = ma
// 	bollUpper = bollMid + float32(PARAM_NUM_SIGMA)*sd
// 	bollLower = bollMid - float32(PARAM_NUM_SIGMA)*sd
// 	//Calculate Parameters of Keltner Channel
// 	kcMid = ema
// 	kcUpper = kcMid + float32(PARAM_NUM_SIGMA)*atr
// 	kcLower = kcMid - float32(PARAM_NUM_SIGMA)*atr
// 	// fmt.Printf("Index:%v Symbol:%s CLOSE: %.2f MA:%.2f SD:%.2f EMA:%.2f ATR:%.2f BMID:%.2f BUP:%.2f BLO:%.2f KMID:%.2f KUP:%.2f KLO:%.2f\r\n", targetIndex, histList[targetIndex].Symbol, histList[targetIndex].MarketClose, ma, sd, ema, atr, bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower)
// 	return bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower
// }
func (sharper *Sharper2) CalcBollKcStat(histList []RowStockHist2, targetIndex int64, crossRange int64, lookbackRange int64) (float32, float32, float32, float32, float32, float32) {
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

	//Calculate MA(MA_RANGE) and SIGMA(MA_RANGE), including the real time stock price
	ma, sd = sharper.CalcMAAndSDStat(histList, targetIndex, PARAM_MA_RANGE)

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
	atr = sharper.CalcATRStat(histList, targetIndex, PARAM_ATR_RANGE)

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

//CalcMAAndSDStat calculates SD (Standard Deviation of Moving Average of Stock Price)
func (sharper *Sharper2) CalcMAAndSDStat(histList []RowStockHist2, targetIndex int64, maRange int64) (float32, float32) {
	var ma float32 //Moving Average
	var sd float32 //Standard Deviation
	for i := targetIndex; i < targetIndex+maRange; i++ {
		ma += (histList[i].GetMarketHigh() + histList[i].GetMarketLow() + 2*histList[i].GetMarketClose()) / 4
	}
	ma /= float32(maRange)
	for i := targetIndex; i < targetIndex+maRange; i++ {
		sdElm := (histList[i].GetMarketHigh()+histList[i].GetMarketLow()+2*histList[i].GetMarketClose())/4 - ma
		sd += float32(math.Pow(float64(sdElm), 2))
	}
	sd = float32(math.Sqrt(float64(sd / float32(maRange))))
	return ma, sd
}

//CalcATRStat calculates ATR (Average True Range)
func (sharper *Sharper2) CalcATRStat(histList []RowStockHist2, targetIndex int64, atrRange int64) float32 {
	var atr float32 //Average True Range
	for i := targetIndex; i < targetIndex+atrRange; i++ {
		trElm1 := float64(histList[i].GetMarketHigh() - histList[i].GetMarketLow())
		trElm2 := math.Abs(float64(histList[i].GetMarketHigh() - histList[i+1].GetMarketClose()))
		trElm3 := math.Abs(float64(histList[i].GetMarketLow() - histList[i+1].GetMarketClose()))
		tr := float32(math.Max(math.Max(trElm1, trElm2), trElm3))
		atr += tr
	}
	atr = atr / float32(atrRange)
	return atr
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

//CalcAvgAndStdDev returns the average and standard deviation of the [0, count) of the data array
//If the length of data array is 0, return error
func (sharper *Sharper2) CalcAvgAndStdDev(data []float32, count int64) (float32, float32, error) {
	var sum float32
	var average float32
	var sumDevSqr float32
	var stdDev float32
	var length int64
	if int64(len(data)) > count {
		length = int64(len(data))
	} else {
		length = count
	}
	if length <= 0 {
		return average, stdDev, errors.New("data array has 0 length or count is non-positive")
	}
	for _, d := range data[0:length] {
		sum += d
	}
	average = sum / float32(length)
	for _, d := range data[0:length] {
		sumDevSqr += (d - average) * (d - average)
	}
	stdDev = float32(math.Sqrt(float64(sumDevSqr / float32(length))))
	return average, stdDev, nil
}

//CalcBlackScholes returns the expected option price using the Black-Scholes Model
func (sharper *Sharper2) CalcBlackScholes(optionType string, underlyingPrice float32, strikePrice float32, impliedVolatility float32, interestRate float32, dividendYield float32, daysToExpire int) float32 {
	timeToExpire := float32(daysToExpire) / float32(365)
	d1 := (float32(math.Log(float64(underlyingPrice)/float64(strikePrice))) + timeToExpire*(interestRate-dividendYield+impliedVolatility*impliedVolatility/2)) / impliedVolatility / float32(math.Sqrt(float64(timeToExpire)))
	d2 := d1 - impliedVolatility*float32(math.Sqrt(float64(timeToExpire)))
	n1 := sharper.CalcNormDistCDF(d1)                                               // normal distribution of d1
	n2 := sharper.CalcNormDistCDF(d2)                                               // normal distribution of d2
	n3 := sharper.CalcNormDistCDF(-d1)                                              //normal distribution of -d1
	n4 := sharper.CalcNormDistCDF(-d2)                                              //normal distribution of -d2
	e1 := strikePrice * float32(math.Exp(float64(-interestRate*timeToExpire)))      // strikePrice * e^(-interestRate * timeToExpire)
	e2 := underlyingPrice * float32(math.Exp(float64(-dividendYield*timeToExpire))) // underlyingPrice * e^(-dividendYield * timeToExpire)
	callPrice := e2*n1 - e1*n2
	putPrice := e1*n4 - e2*n3
	var price float32
	if optionType == "C" {
		price = callPrice
	} else {
		price = putPrice
	}
	return price
}

//CalcNormDistCDF returns the cumulative normal distribution value of mean 0, standard deviation 1
func (sharper *Sharper2) CalcNormDistCDF(value float32) float32 {
	x := float64(value)
	mean := float64(0)
	sigma := float64(1)
	return float32(0.5 * math.Erfc(-(x-mean)/(sigma*math.Sqrt2)))
}

//CalcBlackScholesTrick
func (sharper *Sharper2) CalcBlackScholesTrick(targetMultiple float32, optionType string,
	currentPrice float32, strikePrice float32, impliedVolatility float32,
	interestRate float32, dividendYield float32, currentDaysToExpire int, futureDaysToExpire int) float32 {
	currentOptionPrice := sharper.CalcBlackScholes(optionType, currentPrice, strikePrice, impliedVolatility, interestRate, dividendYield, currentDaysToExpire)
	if currentOptionPrice < 0.1 {
		return -1.0
	}
	var targetOptionPrice float32
	targetOptionPrice = currentOptionPrice * targetMultiple
	var priceLeft float32 = 0.0
	var priceRight float32 = currentPrice * 2
	for priceRight-priceLeft > 0.05 {
		priceMid := priceLeft + (priceRight-priceLeft)/2
		optionPriceMid := sharper.CalcBlackScholes(optionType, priceMid, strikePrice, impliedVolatility, interestRate, dividendYield, futureDaysToExpire)
		if math.Abs(float64(optionPriceMid-targetOptionPrice)) < 0.05 {
			return priceMid
		} else if optionPriceMid > targetOptionPrice {
			if optionType == "C" {
				priceRight = priceMid
			} else {
				priceLeft = priceMid
			}
		} else {
			if optionType == "C" {
				priceLeft = priceMid
			} else {
				priceRight = priceMid
			}
		}
	}
	return -1.0
}
