package main

import (
	c "control"
	"log"
	"net/http"
)

func main() {
	//domain.CustomFun()
	router := c.NewRouter()
	//c.ListKeyspace()
	//c.Dropkeyspace("shivapreals")
	//c.Createkeyspace("shivapreals2")
	//c.DescTables("")
	//d.CustomFun()
	//c.TestRec()
	log.Fatal(http.ListenAndServe(":8081", router))
	//fmt.Println("Shiv")
}
