package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Form struct {
	Id      string `json:"id"`
	Name    string `json:"fullname"`
	Website string `json:"web"`
}

var forms []Form

func (f *Form) IsEmpty() bool {
	return f.Name == ""
}

func main() {
	forms = append(forms, Form{Id: "1", Name: "Pranav", Website: "pranav.something.something"})
	forms = append(forms, Form{Id: "2", Name: "Jain", Website: "jain.something.something"})

	r := mux.NewRouter()
	fmt.Println()
	r.HandleFunc("/", helloHome)
	r.HandleFunc("/all", getAllDetails).Methods("GET")
	r.HandleFunc("/one/{id}", getOneDetail).Methods("GET")
	r.HandleFunc("/create", createDetail).Methods("POST")
	r.HandleFunc("/update/{id}", updateDetail).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteDetails).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}

func helloHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello I am Pranav Jain</h1>"))
}

func getAllDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forms)
}

func getOneDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	parms := mux.Vars(r)
	for _, form := range forms {
		if form.Id == parms["id"] {
			json.NewEncoder(w).Encode(form)
			return
		}
	}
	json.NewEncoder(w).Encode("Not Found")
	return
}

func createDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("No Data sent")
		return
	}
	var form Form
	_ = json.NewDecoder(r.Body).Decode(&form)
	if form.IsEmpty() {
		json.NewEncoder(w).Encode("No Data sent")
		return
	}
	rand.Seed(time.Now().UnixNano())
	form.Id = strconv.Itoa(rand.Intn(100))
	forms = append(forms, form)
	json.NewEncoder(w).Encode(form)
	return
}

func updateDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)

	for index, form := range forms {
		if form.Id == params["id"] {
			forms = append(forms[:index], forms[index+1:]...)
			var form Form
			_ = json.NewDecoder(r.Body).Decode(&form)
			form.Id = params["id"]
			forms = append(forms, form)
			json.NewEncoder(w).Encode(form)
			return
		}
	}
	json.NewEncoder(w).Encode("Id not found")
	return
}

func deleteDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, form := range forms {
		if form.Id == params["id"] {
			forms = append(forms[:index], forms[index+1:]...)
			json.NewEncoder(w).Encode("deleted")
			return
		}
	}
	json.NewEncoder(w).Encode("Id not found")
	return
}
