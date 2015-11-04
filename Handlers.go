package main

import (
  "encoding/json"
  "fmt"
  "net/http"
//  "github.com/gorilla/mux"
  "io/ioutil"
)

/*** THESE ARE THE HANDLERS INFORMATION ***/

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome!")
}

// Return a status for any sendback

func PostHandler(w http.ResponseWriter, r *http.Request) {
  var s Student
  b, _ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(b, &s)
  if err != nil {
    fmt.Printf("Error: json decode")
    fmt.Fprintf(w, "Error: json decode")
  } else {
    fmt.Fprintf(w, "Success")
  }
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
  v := r.URL.Query()
  name := v["name"]
  fmt.Fprintln(w, "Got a student:", name[0])
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
  var y YearQuery
  b, _ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(b, &y)
  if err != nil {
    fmt.Printf("Error: json decode")
    fmt.Fprintf(w, "Error: json decode")
  } else {
    fmt.Fprintln(w, "Year", y.Year)
  }
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "List all")
}
