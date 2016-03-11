package main

import (
	"github.com/junpengxiao/BlockWorldServer/bwRouter"
	"log"
	"net/http"
)

func main() {

	router := bwRouter.NewRouter()
	log.Println("Go Server Starts")
	log.Fatal(http.ListenAndServe(":8080", router))
}
