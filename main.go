package main

import (
	c "control"
	"log"
	"net/http"

	cors "github.com/rs/cors"
)

func main() {
	corsObj := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:1234"},
	})
	//domain.CustomFun()
	router := c.NewRouter()
	//c.ListKeyspace()
	//c.Dropkeyspace("shivapreals")
	//c.Createkeyspace("shivapreals2")
	//c.DescTables("")
	//d.CustomFun()
	//c.TestRec()
	handler := corsObj.Handler(router)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + "8081",
	}

	log.Fatal(srv.ListenAndServe())
	//fmt.Println("Shiv")
}
