package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)


func newRouter() *mux.Router { 
    r := mux.NewRouter()
    r.HandleFunc("/hello", handler).Methods("GET")
    return r
}

func main() {
    // Declaring a new router
    r := newRouter()

    r.HandleFunc("/", handler)

    http.ListenAndServe(":8080", nil)
}


func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}


