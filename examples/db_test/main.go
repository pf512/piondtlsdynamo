package main

import (
	"github.com/pf512/piondtlsdynamo/db"
)
func main() {

	sessiondb.InitDynamo()
	sessiondb.ListAllTables()
	//StoreSession()
	//sessiondb.RetrieveSession()
}

