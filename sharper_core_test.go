package main

import (
	"fmt"
	"testing"
)

func TestSharperCalcAvgAndStdDev(t *testing.T) {
	sharper := new(Sharper2)
	data := []float32{10.2, 9.2, 8.2, 7.2, 6.2}
	ave, sd, baseErr := sharper.CalcAvgAndStdDev(data, 5)
	if baseErr != nil {
		t.Error(baseErr.Error())
	}
	if fmt.Sprintf("%.2f", ave) != "8.20" && fmt.Sprintf("%.2f", sd) != "1.41" {
		t.Fatal("Test Failed: CalcAvgAndStdDev")
	}
	if fmt.Sprintf("%.2f", sharper.CalcNormDistCDF(0.3037619968)) != "0.62" {
		t.Fatal("Test Failed: CalcNormDistCDF")
	}
	c, p := sharper.CalcBlackScholes(36.07, 35, 0.4825, 0.01, 0.0, 26)
	if fmt.Sprintf("%.9f", c) != "2.423076630" || fmt.Sprintf("%.9f", p) != "1.328153610" {
		t.Fatal("Test Failed: CalcBlackScholes")
	}
}
