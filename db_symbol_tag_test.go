package main

import "testing"

func TestDbSymbolTag(t *testing.T) {
	dao := new(DaoSymbolTag)
	if dbDropErr := dao.DropTableIfExist(); dbDropErr != nil {
		t.Error(dbDropErr.Error())
	}
	if dbCreateErr := dao.CreateTableIfNotExist(); dbCreateErr != nil {
		t.Error(dbCreateErr.Error())
	}
	row := RowSymbolTag{
		Symbol: "TEST1",
		Tag:    "TEST2",
	}
	if dbInsertErr := dao.InsertRow(row); dbInsertErr != nil {
		t.Error(dbInsertErr.Error())
	}
	symbols, dbSelectSymbolErr := dao.SelectSymbolByTag(row)
	if dbSelectSymbolErr != nil {
		t.Error(dbSelectSymbolErr.Error())
	}
	if len(symbols) != 1 {
		t.Error("Error: # of rows not matched.")
	}
	if symbols[0].Symbol != "TEST1" {
		t.Errorf("Error: Symbol [0] is not TEST1. It is %s", symbols[0].Symbol)
	}
	tags, dbSelectTagErr := dao.SelectTagBySymbol(row)
	if dbSelectTagErr != nil {
		t.Error(dbSelectTagErr.Error())
	}
	if len(tags) != 1 {
		t.Error("Error: # of rows not matched.")
	}
	if tags[0].Tag != "TEST2" {
		t.Errorf("Error: Tag [0] is not TEST2. It is %s", tags[0].Tag)
	}
	if dbDeleteErr := dao.DeleteRow(row); dbDeleteErr != nil {
		t.Error(dbDeleteErr.Error())
	}
}
