package main

import (
	"fmt"
	"strconv"
	"time"
)

func GetTimeInYYYYMMDD() int {
	i, _ := strconv.Atoi(time.Now().Format("20060102"))
	return i
}
func GetTimeInYYYYMMDD64() int64 {
	i, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 64)
	return i
}
func ConvertTimeInYYYYMMDD(str string) int {
	t, _ := time.Parse("2006-01-02", str)
	i, _ := strconv.Atoi(t.Format("20060102"))
	return i
}

func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetTimeInt() int {
	res, err := strconv.Atoi(time.Now().Format("1504"))
	if err != nil {
		panic(err)
	}
	return res
}

func ConvertTimeInUnix(dateInInt int) int64 {
	t, _ := time.Parse("20060102", strconv.Itoa(dateInInt))
	return t.Unix()
}

func ConvertTime64InUnix(dateInInt64 int64) int64 {
	t, _ := time.Parse("20060102", strconv.FormatInt(dateInInt64, 10))
	return t.Unix()
}

func ConvertUnixTimeInYYYYMMDD(unix int64) int64 {
	t := time.Unix(unix, 0)
	i, _ := strconv.ParseInt(t.Format("20060102"), 10, 64)
	return i
}

func ConvertUTCUnixTimeInYYYYMMDD(unix int64) int64 {
	t := time.Unix(unix, 0)
	loc, err := time.LoadLocation("")
	if err == nil {
		t = t.In(loc)
	}
	i, _ := strconv.ParseInt(t.Format("20060102"), 10, 64)
	return i
}

func MaxInt(i1 int, i2 int) int {
	if i1 < i2 {
		return i2
	} else {
		return i1
	}
}

func MinInt(i1 int, i2 int) int {
	if i1 < i2 {
		return i1
	} else {
		return i2
	}
}

func MaxFloat32(f1 float32, f2 float32) float32 {
	if f1 > f2 {
		return f1
	} else {
		return f2
	}
}

func MinFloat32(f1 float32, f2 float32) float32 {
	if f1 < f2 {
		return f1
	} else {
		return f2
	}
}

func AverageInt(array []int) int {
	if len(array) == 0 {
		return 0
	}
	sum := 0
	for _, num := range array {
		sum += num
	}
	return int(sum / len(array))
}

func AverageInt64(array []int64) int64 {
	if len(array) == 0 {
		return 0
	}
	sum := int64(0)
	for _, num := range array {
		sum += num
	}
	return int64(sum / int64(len(array)))
}

func PrintMsgInConsole(msgType string, logType string, logContent string) {
	fmt.Println(msgType, "[", logType, "] : [", logContent, "]")
}

//ConditionalOperation simulates the conditional operator in Java or C
//Param: statement - the statement of true or false
//Param: trueValue - the returned value if statement is true
//Param: falseValue - the returned value if statement is false
func ConditionalOperation(statement bool, trueValue interface{}, falseValue interface{}) interface{} {
	if statement {
		return trueValue
	} else {
		return falseValue
	}
}
