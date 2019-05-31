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
