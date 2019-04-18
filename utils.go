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
