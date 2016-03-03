package hello

import (
	"bwRouter"
	"net/http"
)

func init() {

	router := bwRouter.NewRouter()
	http.Handle("/", router)
	//	log.Fatal(http.ListenAndServe(":8080", router))
}
