package main

import (
  "encoding/json"
  "fmt"
  "net/http"
//  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
//  "gopkg.in/mgo.v2/bson"
  "io/ioutil"
)

/*** THESE ARE THE HANDLERS INFORMATION ***/

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome!")
}

// Return a status for any sendback

func PostHandler(w http.ResponseWriter, r *http.Request) {
  sess, err0 := mgo.Dial("localhost")
  if err0 != nil {
    fmt.Fprintf(w, "Error: db dial")
    return
  }
  defer sess.Close()

  var s Student
  b, _ := ioutil.ReadAll(r.Body)
  err1 := json.Unmarshal(b, &s)
  if err1 != nil {
    fmt.Fprintf(w, "Error: json decode")
    return
  }
  students := sess.DB("demo1").C("students")
  err2 := students.Insert(s)
  if err2 != nil {
    fmt.Fprintf(w, "Error: db insert")
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
    fmt.Fprintf(w, "Error: json decode")
  } else {
    fmt.Fprintln(w, "Year", y.Year)
  }
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
  sess, err0 := mgo.Dial("localhost")
  if err0 != nil {
    fmt.Fprintf(w, "Error: db dial")
    return
  }
  defer sess.Close()
  students := sess.DB("demo1").C("students")

  iter := students.Find(nil).Iter()
  var s Student
  for iter.Next(s) {
    fmt.Printf("Student: %v\n", s)
  }
}
