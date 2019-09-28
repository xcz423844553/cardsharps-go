package main

import "testing"

func TestDbSymbol(t *testing.T) {
	dao := new(DaoSymbol)
	if dbDropErr := dao.DropTableIfExist(); dbDropErr != nil {
		t.Error(dbDropErr.Error())
	}
	if dbCreateErr := dao.CreateTableIfNotExist(); dbCreateErr != nil {
		t.Error(dbCreateErr.Error())
	}
	row := RowSymbol{
		Symbol: "TEST1",
	}
	if dbInsertOrUpdateErr := dao.InsertOrUpdateRow(row); dbInsertOrUpdateErr != nil {
		t.Error(dbInsertOrUpdateErr.Error())
	}
	row.Symbol = "TEST2"
	if dbInsertErr := dao.InsertRow(row); dbInsertErr != nil {
		t.Error(dbInsertErr.Error())
	}
	if dbUpdateErr := dao.UpdateRow(row); dbUpdateErr != nil {
		t.Error(dbUpdateErr.Error())
	}
	if dbDeleteErr := dao.DeleteRow(row); dbDeleteErr != nil {
		t.Error(dbDeleteErr.Error())
	}
	rows, dbSelectAllErr := dao.SelectSymbolAll()
	if dbSelectAllErr != nil {
		t.Error(dbSelectAllErr.Error())
	}
	if len(rows) != 1 {
		t.Error("Error: # of rows not matched.")
	}
	if rows[0].Symbol != "TEST1" {
		t.Errorf("Error: Symbol [0] is not TEST1. It is %s", rows[0].Symbol)
	}
}
