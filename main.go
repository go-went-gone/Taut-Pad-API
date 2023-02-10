package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"log"
	"github.com/gorilla/mux"
)

type Note struct{
	Title string `json:"title"`
	Description string `json:"description"`
}

var noteStore = make(map[string]Note)

var id int = 0

func tautPadHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case "POST":
		var note Note
		//Decode the incoming Note Json
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil{
			panic(err)
		}
		id++
		k := strconv.Itoa(id)
		noteStore[k] = note
		//Convert to json type
		j,err := json.Marshal(note)
		if err != nil{
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(j)
	case "GET":
		var notes []Note
		for _,v := range noteStore{
			notes = append(notes, v)
		}
		w.Header().Set("Content-Type", "application/json")
		j,err := json.Marshal(notes)
		if err != nil{
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}


func getDetailAndDeleteTautPadHandler(w http.ResponseWriter, r * http.Request){
	switch r.Method{
	case "PUT":
		var err error
		vars := mux.Vars(r)
		k := vars["id"]
		var noteToUpdate Note
		err = json.NewDecoder(r.Body).Decode(&noteToUpdate)
		if err != nil{
			panic(err)
		}
		if _,ok := noteStore[k]; ok{
			delete(noteStore, k)
			noteStore[k]=noteToUpdate
		}else{
			log.Printf("Couldn't find the resource key %v to delete", k)
		}
	case "DELETE":
		vars := mux.Vars(r)
		k := vars["id"]
		if _, ok := noteStore[k]; ok{
			delete(noteStore, k)
		}else{
			log.Printf("Couldn't find the resource key %v to delete", k)
		}
	default: 
		w.WriteHeader(http.StatusNoContent)
	}
	
	



}

func Logging(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start := time.Now()
		log.Printf("Starting Server...")
		if  r.Method  == "PUT"{
			log.Printf("Started %v %v", r.Method, r.URL.Path)
		}
		h.ServeHTTP(w, r)
		log.Printf("Completed %v in %v", r.URL.Path, time.Since(start))
	})
}

func main(){
	router := mux.NewRouter().StrictSlash(false)
	
	router.HandleFunc("/notes", tautPadHandler).Methods("GET","POST")
	router.HandleFunc("/notes/{id}", getDetailAndDeleteTautPadHandler).Methods("PUT","DELETE")

	server := &http.Server{
		Addr : "0.0.0.0:8080",
		Handler : router,
	}

	server.ListenAndServe()
}

