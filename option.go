package main

import (
	"math"
)

type OptionFilter struct {
	MaxOptionPercent  float32 `json:"maxOptionPercent"`
	MinOptionPercent  float32 `json:"minOptionPercent"`
	MaxOpenInterest   int64   `json:"maxOpenInterest"`
	MinOpenInterest   int64   `json:"minOpenInterest"`
	MaxVolume         int64   `json:"maxVolume"`
	MinVolume         int64   `json:"minVolume"`
	MaxExpirationDate int64   `json:"maxExpirationDate"`
	MinExpirationDate int64   `json:"minExpirationDate"`
}

//func IsInOptionFilter

func NewOptionFilter(maxOptionPercent float32, minOptionPercent float32,
	maxOpenInterest int64, minOpenInterest int64, maxVolume int64, minVolume int64,
	maxExpirationDate int64, minExpirationDate int64) OptionFilter {
	f := new(OptionFilter)
	if maxOptionPercent == 0.0 {
		f.MaxOptionPercent = 1.0
	} else {
		f.MaxOptionPercent = maxOptionPercent
	}
	if minOptionPercent == 0.0 {
		f.MinOptionPercent = -1.0
	} else {
		f.MinOptionPercent = minOptionPercent
	}
	if maxOpenInterest == 0 {
		f.MaxOpenInterest = math.MaxInt64
	} else {
		f.MaxOpenInterest = maxOpenInterest
	}
	if minOpenInterest == 0 {
		f.MinOpenInterest = math.MinInt64
	} else {
		f.MinOpenInterest = minOpenInterest
	}
	if maxVolume == 0 {
		f.MaxVolume = math.MaxInt64
	} else {
		f.MaxVolume = maxVolume
	}
	if minVolume == 0 {
		f.MinVolume = math.MinInt64
	} else {
		f.MinVolume = minVolume
	}
	if maxExpirationDate == 0 {
		f.MaxExpirationDate = math.MaxInt64
	} else {
		f.MaxExpirationDate = maxExpirationDate
	}
	if minExpirationDate == 0 {
		f.MinExpirationDate = math.MinInt64
	} else {
		f.MinExpirationDate = minExpirationDate
	}
	return *f
}
