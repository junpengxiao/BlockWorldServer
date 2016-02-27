package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//"github.com/gorilla/mux"
)

var usage = `Welcome!
/view :  To test/view world with JSON form, post it to /view
	eg: curl -H "Content-T{"loc": [-0.1667, 0.1, -0.3333], "id": 1}], "version":1, "error":"Null"}' http://localhost:8080/view
/query:  To get the prediction of the world with JSON form, post it to /query which will return a predicted world with JSON form
/upload: return a webpage. To test server locally, use this page to upload a JSON file
`

func obtainJSON(w http.ResponseWriter, r *http.Request) (Data, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	return data, err
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, usage)
}

func View(w http.ResponseWriter, r *http.Request) {
	data, err := obtainJSON(w, r)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, data)
}

var uploadpage = `<html>
	<body>
		<form enctype="multipart/form-data" action="/viewupload" method="post">
			<input type="file" name="data">
			<input type="submit">
		</form>
	</body>
</html>
`

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, uploadpage)
}

func ViewUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	/*
		var data Data
		if err = json.Unmarshal(r.Form["data"], &data); err!=nil {
			panic(err)
		}
		result := ProcessData(data)
		//remember to return result
	*/
	file, _, err := r.FormFile("data")
	if err != nil {
		fmt.Fprint(w, err)
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	var data Data
	err = json.Unmarshal(body, &data)
	fmt.Fprint(w, data)
}

func Query(w http.ResponseWriter, r *http.Request) {
	data, err := obtainJSON(w, r)
	if err != nil {
		return
	}
	//result := ProcessData(data)
	result := data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}