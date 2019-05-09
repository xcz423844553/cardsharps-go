package main

import (
	"strconv"
	"time"
)

func GetTimeInYYYYMMDD() int {
	i, _ := strconv.Atoi(time.Now().Format("20060102"))
	return i
}

func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
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
