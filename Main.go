package main

// TODO: create REAL test.go

import (
  "log"
  "net/http"
  "encoding/json"
  "fmt"
  "gopkg.in/mgo.v2"
//  "gopkg.in/mgo.v2/bson"
  "io/ioutil"
  "time"
  "github.com/gorilla/mux"
)

////////////////////////////
//      GLOBALS
////////////////////////////

// This is the hack that caused us to use a single file
var Database *mgo.Database

////////////////////////////
//      Main.go
////////////////////////////

func main() {
  session, err := mgo.Dial("mongodb://localhost:27017")
  if err != nil {
    fmt.Printf("Error: bad mgo")
    return
  }
  Database = session.DB("hw3")

  router := NewRouter()
  log.Fatal(http.ListenAndServe(":1234", router))
}

////////////////////////////
//      Student.go
////////////////////////////

// TODO: change the string type of netid to bson.ObjectId, then test
type Student struct {
  NetID string `json: "id" bson: "id"`
  Name string `json: "name" bson: "name"`
  Major string `json: "major" bson: "major"`
  Year int `json: "year" bson: "year"`
  Grade int `json: "grade" bson: "grade"`
  Rating string `json: "rating" bson: "rating"`
}

type Students []Student

////////////////////////////
//      YearQuery.go
////////////////////////////

type YearQuery struct {
  Year int `json: "year"`
}

////////////////////////////
//      Routes.go
////////////////////////////

type Route struct {
  Name string
  Method string
  Pattern string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    var handler http.Handler
    handler = route.HandlerFunc
    handler = Logger(handler, route.Name)
    // attach the routes to their handlers
    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)
  }
  return router
}

/*** THIS IS WHERE WE PUT ALL THE ROUTING INFORMATION ***/
/*** COMMUNICATES WITH HANDLERS.GO ***/

var routes = Routes{
  Route {
    "IndexPage",
    "GET",
    "/",
    Index,
  },
  Route {
    "POST_Operation",
    "POST",
    "/Student",
    PostHandler,
  },
  Route {
    "GET_Operation",
    "GET",
    "/Student/getstudent",
    GetHandler,
  },
  Route {
    "DELETE_Operation",
    "DELETE",
    "/Student",
    DeleteHandler,
  },
  Route {
    "UPDATE_Operation",
    "PUT",
    "/Student",
    UpdateHandler,
  },
  Route {
    "LIST_Operation",
    "GET",
    "/Student/listall",
    ListHandler,
  },
}

////////////////////////////
//      Logger.go
////////////////////////////

func Logger(inner http.Handler, name string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    inner.ServeHTTP(w, r)
    log.Printf("%s\t%s\t%s\t%s",
      r.Method,
      r.RequestURI,
      name,
      time.Since(start),
    )
  })
}

////////////////////////////
//      Routes.go
////////////////////////////

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome!\nProject by: Shahan, Edward, Shina, Tim\n")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
  // decode the json
  var s Student
  b, _ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(b, &s)
  if err != nil {
    fmt.Fprintf(w, "Error: json decode\n")
    return
  }
  // store in database
  students := Database.C("students").With(Database.Session.Copy())
  err = students.Insert(s)
  if err != nil {
    fmt.Fprintf(w, "Error: mongodb insert\n")
  } else {
    fmt.Fprintf(w, "Success\n")
  }
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
  v := r.URL.Query()
  name := v["name"]
  // change this to grab database entry
  fmt.Fprintln(w, "Got a student:", name[0])
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
  var y YearQuery
  b, _ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(b, &y)
  // once has the name, delete it
  if err != nil {
    fmt.Fprintf(w, "Error: json decode")
  } else {
    fmt.Fprintln(w, "Year", y.Year)
  }
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
  // add stuff in here
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
  students := Database.C("students").With(Database.Session.Copy())
  var s Student
  iter := students.Find(nil).Iter()
  for iter.Next(&s) {
    fmt.Fprintf(w, "Student: %v\n", s)
  }
}
