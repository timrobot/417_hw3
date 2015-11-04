package main

// TODO: create REAL test.go

import (
  "log"
  "net/http"
  "encoding/json"
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "io/ioutil"
  "time"
  "github.com/gorilla/mux"
  "os"
)

////////////////////////////
//      Main.go
////////////////////////////

func main() {
  router := NewRouter()
  log.Fatal(http.ListenAndServe(":1234", router))
}

////////////////////////////
//      Student.go
////////////////////////////

// TODO: change the string type of netid to bson.ObjectId, then test

type Student struct {
    NetID  string `json:"id" bson:"_id"`  
    Name  string  `json:"name" bson:"name"`
    Major  string `json:"major" bson:"major"`
    Year  int `json:"year" bson:"year"`
    Grade  int  `json:"grade" bson:"grade"`
    Rating  string  `json:"rating" bson:"rating"`
}

type Students []Student

////////////////////////////
//      YearQuery.go
////////////////////////////

type YearQuery struct {
  Year int `json: "year" bson: "year"`
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
//      Handlers.go
////////////////////////////

func InitSession() *mgo.Session {
  s, err := mgo.Dial("mongodb://localhost:27017")
  if err != nil{
    fmt.Printf("Error: db init\n")
    os.Exit(1)
  } 
  s.SetSafe(&mgo.Safe{})
  return s
}

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
  students := InitSession().DB("hw3").C("students")
  err = students.Insert(s)
  if err != nil {
    fmt.Fprintf(w, "Error: mongodb insert\n")
  } else {
    fmt.Fprintf(w, "Success\n")
  }
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
  v := r.URL.Query()
  name := v["name"][0]
  students := InitSession().DB("hw3").C("students")
  var s *Student
  err := students.Find(bson.M{"name": name}).One(&s)
  if err != nil {
    fmt.Fprintf(w, "Error: db lookup\n")
    return
  }
  student, _ := json.Marshal(s)
  fmt.Fprintf(w, "Student: %s\n", student)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
  var y YearQuery
  b, _ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(b, &y)
  if err != nil {
    fmt.Fprintf(w, "Error: json decode\n")
    return
  }
  students := InitSession().DB("hw3").C("students")
  var s []Student
  err = students.Find(nil).All(&s)
  if err != nil {
    fmt.Fprintf(w, "Error: delete\n")
    return
  }
  for i := 0; i < len(s); i++ {
    if s[i].Year < y.Year {
      students.RemoveId(s[i].NetID)
    }
  }
  fmt.Fprintf(w, "Success\n")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
  students := InitSession().DB("hw3").C("students")
  var s []Student
  average := 0
  rating := "E"
  err := students.Find(nil).All(&s)
  if err != nil {
    fmt.Fprintf(w, "Error: update\n")
    return
  }
  for i := 0; i < len(s); i++ {
    average += s[i].Grade
  }
  average = average / len(s)
  for i := 0; i < len(s); i++ {
    if s[i].Grade > average + 10 {
      rating = "A"
    } else if average - 10 < s[i].Grade && s[i].Grade <= average + 10 {
      rating = "B"
    } else if average - 20 < s[i].Grade && s[i].Grade <= average - 10 {
      rating = "C"
    } else {
      rating = "D"
    }
    students.UpdateId(s[i].NetID, bson.M{"$set": bson.M{"rating": rating}})  
  }
  fmt.Fprintf(w, "Success\n")
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
  students := InitSession().DB("hw3").C("students")
  var s []Student
  err := students.Find(nil).All(&s)
  if err != nil {
    fmt.Fprintf(w, "Error: db list")
    return
  }
  listdata, _ := json.Marshal(s)
  fmt.Fprintf(w, "Students:\n%s\n", listdata)
}
