package main

import (
	"github.com/pion/dtls/v2/db"
)
func main() {

	sessiondb.InitDynamo()
	sessiondb.ListAllTables()
	//StoreSession()
	//sessiondb.RetrieveSession()
}

