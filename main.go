package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func getGreeting() string {
	return "Hello, Kontur!"
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, getGreeting())
}

func main() {
	//http.ListenAndServe(":8080", nil)
	repo := New()
	srv := NewServer(repo)
	router := mux.NewRouter()

	router.HandleFunc("/hello", hello)
	router.HandleFunc("/chartas/", srv.NewCharta).Methods("POST")
	router.HandleFunc("/chartas/{id}", srv.EditCharta).Methods("POST")
	router.HandleFunc("/chartas/{id}", srv.GetCharta).Methods("GET")
	router.HandleFunc("/chartas/{id}", srv.DeleteCharta).Methods("DELETE")

	port := "8000" //localhost

	//fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}

}
